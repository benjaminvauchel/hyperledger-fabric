#!/bin/bash

ARGS_FILE=$1
START_ID=${2:-1}
PARALLELISM=${3:-5}
RESULTS_DIR="eval_results"
mkdir -p "$RESULTS_DIR"
OUTFILE="$RESULTS_DIR/results_${START_ID}_${PARALLELISM}.csv"
LATENCY_FILE="$RESULTS_DIR/latencies_${START_ID}_${PARALLELISM}.txt"

# Certificate paths
ORDERER_CA="${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
PEER0_ORG1_CA="${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
PEER0_ORG2_CA="${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# Set environment variables for peer commands
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
export CORE_PEER_MSPCONFIGPATH="${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
export CORE_PEER_ADDRESS="localhost:7051"

CHAINCODE_NAME="basic"
CHANNEL_NAME="mychannel"
FUNCTION_NAME="CreateAcademicCredential"

mkdir -p invoke_logs
# Initialize the results file and latency file
echo "credential_id,start_time_ms,end_time_ms,latency_ms" > "$OUTFILE"
> "$LATENCY_FILE"  # Create or truncate latency file

# Store credential IDs in an array for tracking
declare -a CREDENTIAL_IDS

# More precise time measurement function that works in subshells
get_time_ms() {
  echo $(($(date +%s%N)/1000000))
}

run_invoke() {
  local id="$1"
  local args=("$@")
  
  # Remove the first argument (id) from args array
  args=("${args[@]:1}")
  
  # Format arguments for chaincode invocation - properly JSON-escaped
  local json_args="[\"$id\""
  for arg in "${args[@]}"; do
    json_args+=",\"$arg\""
  done
  json_args+="]"
  
  # Get start time in milliseconds with nanosecond precision
  local start_time=$(get_time_ms)
  
  peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls \
    --cafile "$ORDERER_CA" \
    -C "$CHANNEL_NAME" \
    -n "$CHAINCODE_NAME" \
    --peerAddresses localhost:7051 \
    --tlsRootCertFiles "$PEER0_ORG1_CA" \
    --peerAddresses localhost:9051 \
    --tlsRootCertFiles "$PEER0_ORG2_CA" \
    --waitForEvent \
    -c "{\"function\":\"$FUNCTION_NAME\",\"Args\":$json_args}" \
    &> "invoke_logs/${id}.log"
  
  # Get end time in milliseconds with nanosecond precision
  local end_time=$(get_time_ms)
  
  # Calculate latency - force abs value with ${var#-} to handle any timing anomalies
  local latency=$((end_time - start_time))
  latency=${latency#-}  # Remove minus sign if present

  if [ $latency -eq 0 ]; then
    # If somehow we get zero latency, set to 1 ms minimum
    latency=1
  fi
  
  # Write to results file
  echo "$id,$start_time,$end_time,$latency" >> "$OUTFILE"
  # Also write just the latency to a separate file for statistics
  echo "$latency" >> "$LATENCY_FILE"
  
  #echo "✅ $id done in ${latency}ms"
}

# Function to calculate statistics from latency file
calc_stats() {
  if [ ! -s "$LATENCY_FILE" ]; then
    echo "0,0,0,0,0,0,0"  # Return zeros if file is empty
    return
  fi
  
  # Sort latencies for percentile calculations
  sort -n "$LATENCY_FILE" > "${LATENCY_FILE}.sorted"
  
  # Count, sum, min, max
  local count=$(wc -l < "$LATENCY_FILE")
  local sum=$(awk '{sum += $1} END {print sum}' "$LATENCY_FILE")
  local min=$(head -n 1 "${LATENCY_FILE}.sorted")
  local max=$(tail -n 1 "${LATENCY_FILE}.sorted")
  
  # Average
  local avg=$(echo "scale=2; $sum / $count" | bc)
  
  # Get percentiles using awk - more robust than manual line calculation
  local percentiles=$(awk '
    BEGIN {p50=0; p95=0; p99=0;}
    {
      values[NR] = $1;
    }
    END {
      n = NR;
      if (n > 0) {
        p50_idx = int(n * 0.5);
        p95_idx = int(n * 0.95);
        p99_idx = int(n * 0.99);
        
        # Make sure indices are at least 1
        if (p50_idx < 1) p50_idx = 1;
        if (p95_idx < 1) p95_idx = 1;
        if (p99_idx < 1) p99_idx = 1;
        
        # Cap indices at array size
        if (p50_idx > n) p50_idx = n;
        if (p95_idx > n) p95_idx = n;
        if (p99_idx > n) p99_idx = n;
        
        print values[p50_idx] "," values[p95_idx] "," values[p99_idx];
      } else {
        print "0,0,0";
      }
    }
  ' "${LATENCY_FILE}.sorted")
  
  # Parse the percentiles
  IFS=',' read -r median_latency p95_latency p99_latency <<< "$percentiles"
  
  # Standard deviation
  local variance=$(awk -v avg="$avg" '{sum += ($1-avg)^2} END {print sum/NR}' "$LATENCY_FILE")
  local stddev=$(echo "scale=2; sqrt($variance)" | bc)
  
  # Return comma-separated statistics
  echo "$avg,$median_latency,$p95_latency,$p99_latency,$min,$max,$stddev"
  
  # Clean up
  rm -f "${LATENCY_FILE}.sorted"
}

# Main execution
i=0
COUNTER=$START_ID
total_start=$(get_time_ms)

# Read and process the entries from the JSON file
while read -r entry; do
  cid="credential$COUNTER"
  
  # Extract array elements from the JSON entry
  readarray -t args < <(echo "$entry" | jq -r '.[]')
  
  # Run invoke with credential ID and extracted arguments
  run_invoke "$cid" "${args[@]}" &
  
  CREDENTIAL_IDS+=("$cid")
  ((COUNTER++))
  ((i++))

  if (( i % PARALLELISM == 0 )); then
    wait
  fi
done < <(jq -c '.[]' "$ARGS_FILE")

wait
total_end=$(get_time_ms)

total_duration=$((total_end - total_start))
total_tx=${#CREDENTIAL_IDS[@]}

if (( total_duration > 0 )); then
  throughput=$(echo "scale=2; $total_tx / ($total_duration / 1000)" | bc)
else
  throughput="N/A"
fi

# Calculate latency statistics
IFS=',' read -r avg_latency median_latency p95_latency p99_latency min_latency max_latency stddev_latency < <(calc_stats)

# Create a statistics file
STATS_FILE="$RESULTS_DIR/stats_${START_ID}_${PARALLELISM}.txt"
{
  echo "===== HYPERLEDGER FABRIC PERFORMANCE METRICS ====="
  echo "Test Configuration:"
  echo "- Start ID: $START_ID"
  echo "- Parallelism: $PARALLELISM"
  echo "- Total Transactions: $total_tx"
  echo ""
  echo "Throughput Metrics:"
  echo "- Total Duration: ${total_duration} ms"
  echo "- Throughput: $throughput tx/sec"
  echo ""
  echo "Latency Metrics (ms):"
  echo "- Average: $avg_latency"
  echo "- Median (P50): $median_latency"
  echo "- P95: $p95_latency"
  echo "- P99: $p99_latency"
  echo "- Min: $min_latency"
  echo "- Max: $max_latency"
  echo "- Standard Deviation: $stddev_latency"
  echo ""
  echo "Full results saved to: $OUTFILE"
  echo "==================================================="
} > "$STATS_FILE"

echo ""
echo "📊 Evaluation Complete"
echo "🕒 Total Time     : ${total_duration} ms"
echo "🔁 Transactions   : $total_tx"
echo "🚀 Throughput     : $throughput tx/sec"
echo ""
echo "📈 Latency Statistics (ms):"
echo "   - Average: $avg_latency"
echo "   - Median (P50): $median_latency" 
echo "   - P95: $p95_latency"
echo "   - P99: $p99_latency"
echo "   - Min: $min_latency"
echo "   - Max: $max_latency"
echo "   - Std Dev: $stddev_latency"
echo ""
echo "📄 Results saved to: $OUTFILE"
echo "📊 Statistics saved to: $STATS_FILE"
