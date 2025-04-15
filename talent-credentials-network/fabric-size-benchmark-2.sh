#!/bin/bash

# Default settings
ITERATIONS=${1:-100}          # Number of transactions per size
PARALLELISM=${2:-5}           # Fixed parallelism level
START_ID=${3:-70000}          # Starting credential ID
RESULTS_DIR="size_benchmarks" # Results directory

# Create directories
mkdir -p "$RESULTS_DIR"
mkdir -p "test_data"
mkdir -p "invoke_logs"

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

# Size configurations - sizes in bytes
SIZES=(64 128 512 32)

# Function to generate random string of specified length
generate_random_string() {
  local length=$1
  if [ $length -le 0 ]; then
    echo ""
    return
  fi
  cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w $length | head -n 1
}

# Function to get current time in milliseconds
get_time_ms() {
  echo $(($(date +%s%N)/1000000))
}

# Function to calculate exact payload size
calculate_payload_size() {
  local json_args="$1"
  echo -n "$json_args" | wc -c
}

# Function to generate test data for each size
generate_test_data() {
  local target_size=$1
  local output_file="test_data/args_${target_size}B.json"
  
  echo "[" > "$output_file"
  
  for ((i=1; i<=$ITERATIONS; i++)); do
    # Calculate the credential ID string length
    local cred_id="cred$((START_ID + i - 1))"
    local cred_id_len=${#cred_id}
    
    # The full JSON structure will be: ["cred80100","a","b","c","d","e","f"]
    # Let's calculate the JSON overhead precisely:
    # - 2 bytes for outer brackets '[]'
    # - 2 bytes for quotes around credential ID
    # - cred_id_len (9) bytes for the credential ID itself
    # - 18 bytes for all the commas and quotes separating the other fields
    
    # Total JSON overhead
    local json_overhead=$((2 + 2 + cred_id_len + 18))
    
    # Calculate space available for actual field data
    local available_space=$((target_size - json_overhead))
    
    if [ $available_space -lt 0 ]; then
      # Handle too small target size
      echo "  [\"a\",\"b\",\"c\",\"d\",\"e\",\"f\"]" >> "$output_file"
      echo "Warning: Target size ${target_size}B is smaller than minimum required (${json_overhead}B)" >&2
      continue
    fi
    
    # Distribute available space across fields evenly
    local base_size=$((available_space / 6))
    local extra=$((available_space % 6))
    
    local talent_size=$base_size
    local first_size=$base_size
    local last_size=$base_size
    local skills_size=$base_size
    local degree_size=$base_size
    local university_size=$base_size
    
    # Distribute any remainder
    if [ $extra -gt 0 ]; then ((talent_size++)); ((extra--)); fi
    if [ $extra -gt 0 ]; then ((first_size++)); ((extra--)); fi
    if [ $extra -gt 0 ]; then ((last_size++)); ((extra--)); fi
    if [ $extra -gt 0 ]; then ((skills_size++)); ((extra--)); fi
    if [ $extra -gt 0 ]; then ((degree_size++)); ((extra--)); fi
    if [ $extra -gt 0 ]; then ((university_size++)); fi
    
    # Generate data of exact sizes
    local talent_id=$(generate_random_string $talent_size)
    local first_name=$(generate_random_string $first_size)
    local last_name=$(generate_random_string $last_size)
    local skills=$(generate_random_string $skills_size)
    local degree=$(generate_random_string $degree_size)
    local university=$(generate_random_string $university_size)
    
    # Write to the output file
    echo -n "  [\"$talent_id\",\"$first_name\",\"$last_name\",\"$skills\",\"$degree\",\"$university\"]" >> "$output_file"
    
    # Add comma if not the last entry
    if [ $i -lt $ITERATIONS ]; then
      echo "," >> "$output_file"
    else
      echo "" >> "$output_file"
    fi
  done
  
  echo "]" >> "$output_file"
  echo "  Generated ${ITERATIONS} entries for ${target_size}B payload"
}

# Function to run transaction with given arguments and specific size
run_invoke() {
  local size="$1"
  local id="$2"
  shift 2
  local args=("$@")
  
  # Format arguments for chaincode invocation
  local json_args="[\"$id\""
  for arg in "${args[@]}"; do
    json_args+=",\"$arg\""
  done
  json_args+="]"
  
  # Verify the exact size of the payload
  local actual_size=$(calculate_payload_size "$json_args")
  
  # Log the size match if it doesn't match target
  if [ "$actual_size" != "$size" ]; then
    echo "⚠️ Size mismatch for $id: Target=$size, Actual=$actual_size" >&2
    
    # Debug the structure
    if [ -n "${DEBUG:-}" ]; then
      echo "DEBUG: JSON structure: $json_args" >&2
      echo "DEBUG: Character count: $(echo -n "$json_args" | wc -c)" >&2
      echo "DEBUG: Hex dump:" >&2
      echo -n "$json_args" | hexdump -C >&2
    fi
  fi
  
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
  
  # Calculate latency
  local latency=$((end_time - start_time))
  latency=${latency#-}  # Remove minus sign if present (to handle any timing anomalies)

  if [ $latency -eq 0 ]; then
    latency=1  # If somehow we get zero latency, set to 1 ms minimum
  fi
  
  # Write to results file
  echo "$id,$actual_size,$start_time,$end_time,$latency" >> "$OUTFILE"
  # Also write just the latency to a separate file for statistics
  echo "$latency" >> "$LATENCY_FILE"
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
  if [ $count -eq 0 ]; then
    echo "0,0,0,0,0,0,0"  # Return zeros if count is zero
    return
  fi
  
  local sum=$(awk '{sum += $1} END {print sum}' "$LATENCY_FILE")
  local min=$(head -n 1 "${LATENCY_FILE}.sorted")
  local max=$(tail -n 1 "${LATENCY_FILE}.sorted")
  
  # Average - protect against divide by zero
  if [ $count -gt 0 ]; then
    local avg=$(echo "scale=2; $sum / $count" | bc)
  else
    local avg=0
  fi
  
  # Get percentiles using awk with safety checks
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
  
  # Standard deviation with safety check
  if [ $count -gt 1 ]; then
    local variance=$(awk -v avg="$avg" '{sum += ($1-avg)^2} END {print sum/NR}' "$LATENCY_FILE")
    local stddev=$(echo "scale=2; sqrt($variance)" | bc)
  else
    local stddev=0
  fi
  
  # Return comma-separated statistics
  echo "$avg,$median_latency,$p95_latency,$p99_latency,$min,$max,$stddev"
  
  # Clean up
  rm -f "${LATENCY_FILE}.sorted"
}

# Print header
echo "=========================================="
echo "Hyperledger Fabric Size Benchmark"
echo "Testing payload sizes: ${SIZES[*]} bytes"
echo "Starting credential ID: $START_ID"
echo "Parallelism: $PARALLELISM"
echo "Iterations per size: $ITERATIONS"
echo "=========================================="

# Generate test data for each size
echo "Generating test data for different payload sizes..."
for size in "${SIZES[@]}"; do
  echo "- Generating ${size}B payload data..."
  generate_test_data $size
done
echo "Test data generation complete."

# Create summary file for all results
SUMMARY_FILE="$RESULTS_DIR/size_benchmark_summary.csv"
echo "payload_size,avg_latency,median_latency,p95_latency,p99_latency,min_latency,max_latency,stddev_latency,throughput" > "$SUMMARY_FILE"

# Run benchmarks for each size
echo "Starting benchmarks for different payload sizes..."
current_id=$START_ID

for size in "${SIZES[@]}"; do
  echo "Running benchmark for ${size}B payloads..."
  
  # Setup files for this size
  OUTFILE="$RESULTS_DIR/results_${size}B.csv"
  LATENCY_FILE="$RESULTS_DIR/latencies_${size}B.txt"
  STATS_FILE="$RESULTS_DIR/stats_${size}B.txt"
  
  # Initialize results files
  echo "credential_id,actual_size,start_time_ms,end_time_ms,latency_ms" > "$OUTFILE"
  > "$LATENCY_FILE" # Create or truncate latency file
  
  ARGS_FILE="test_data/args_${size}B.json"
  
  # Start timing for throughput calculation
  total_start=$(get_time_ms)
  
  # Process the test data
  i=0
  while read -r entry; do
    cid="cred$current_id"
    
    # Extract array elements from the JSON entry
    readarray -t args < <(echo "$entry" | jq -r '.[]')
    
    # Run invoke with credential ID and extracted arguments
    run_invoke "$size" "$cid" "${args[@]}" &
    
    ((current_id++))
    ((i++))
    
    if (( i % PARALLELISM == 0 )); then
      wait
    fi
  done < <(jq -c '.[]' "$ARGS_FILE")
  
  wait
  total_end=$(get_time_ms)
  
  # Calculate metrics
  total_duration=$((total_end - total_start))
  total_tx=$ITERATIONS
  
  if (( total_duration > 0 && total_tx > 0 )); then
    throughput=$(echo "scale=2; $total_tx / ($total_duration / 1000)" | bc)
  else
    throughput="0.00"
  fi
  
  # Calculate latency statistics
  IFS=',' read -r avg_latency median_latency p95_latency p99_latency min_latency max_latency stddev_latency < <(calc_stats)
  
  # Add to summary
  echo "${size},${avg_latency},${median_latency},${p95_latency},${p99_latency},${min_latency},${max_latency},${stddev_latency},${throughput}" >> "$SUMMARY_FILE"
  
  # Create a statistics file for this size
  {
    echo "===== HYPERLEDGER FABRIC SIZE BENCHMARK: ${size}B ====="
    echo "Test Configuration:"
    echo "- Payload Size: ${size} bytes"
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
  
  echo "✅ ${size}B benchmark complete"
  echo ""
done

# Create visualization data (for easy plotting)
echo "size,metric,value" > "$RESULTS_DIR/plot_data.csv"
while IFS=',' read -r size avg med p95 p99 min max stddev tput; do
  # Skip header
  if [ "$size" != "payload_size" ]; then
    echo "$size,avg_latency,$avg" >> "$RESULTS_DIR/plot_data.csv"
    echo "$size,median_latency,$med" >> "$RESULTS_DIR/plot_data.csv"
    echo "$size,p95_latency,$p95" >> "$RESULTS_DIR/plot_data.csv"
    echo "$size,p99_latency,$p99" >> "$RESULTS_DIR/plot_data.csv"
    echo "$size,throughput,$tput" >> "$RESULTS_DIR/plot_data.csv"
  fi
done < "$SUMMARY_FILE"

echo ""
echo "📊 Size Benchmark Complete"
echo "📄 Summary saved to: $SUMMARY_FILE"
echo "📈 Plot data available in: $RESULTS_DIR/plot_data.csv"
echo ""
echo "Size benchmarks:"
cat "$SUMMARY_FILE"
