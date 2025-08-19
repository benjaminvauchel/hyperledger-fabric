# Blockchain-Based System for Secure Talent Credential Verification

A decentralized, tamper-resistant solution for verifying academic and professional credentials using Hyperledger Fabric blockchain technology.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#theoritical-usage)
- [API Documentation](#api-documentation)
- [Smart Contracts](#smart-contracts)
- [Performance Results](#performance-results)
- [License](#license)

## Overview

This project addresses the critical challenge of verifying academic and professional credentials by leveraging blockchain technology. Traditional credential verification systems suffer from inefficiencies, forgery risks, and lack of standardized verification mechanisms. Our solution provides a secure, decentralized platform where talents, institutions, and companies can interact to create, verify, and manage credentials with complete transparency and security.

## Features

- **Decentralized Verification**: Permissioned blockchain network ensuring data integrity
- **Multi-Organization Support**: Three distinct roles - Talents, Institutions, and Companies
- **Smart Contract Automation**: Automated credential lifecycle management
- **REST API Integration**: Seamless backend-frontend communication
- **Web Interface**: User-friendly frontend for credential management
- **Tamper-Resistant Storage**: Immutable credential records on the blockchain
- **Role-Based Access Control**: Granular permissions for different user types

## Architecture

The system is built on Hyperledger Fabric with the following components:

### Network Structure
- **3 Organizations**: Talents, Institutions, Companies
- **Peers**: Multiple peer nodes across organizations
- **Ordering Service**: RAFT-based consensus with fault tolerance
- **Certificate Authorities**: X.509 certificate management for each organization
- **Shared Channel**: All participants interact over a common channel

### Technology Stack
- **Blockchain**: Hyperledger Fabric
- **Smart Contracts**: Go (Golang)
- **Backend API**: Node.js with Hyperledger Fabric Gateway SDK
- **Frontend**: Web application with REST API integration
- **Deployment**: Docker containers
- **Identity Management**: X.509 certificates and MSPs

## Prerequisites

- Docker and Docker Compose
- Node.js (v14 or higher)
- Go (v1.19 or higher)
- Hyperledger Fabric binaries and Docker images
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/benjaminvauchel/hyperledger-fabric.git
   cd hyperledger-fabric
   ```

2. **Set up Hyperledger Fabric**
   ```bash
   # Download Fabric binaries and Docker images
   ```

3. **Deploy the network**
   ```bash
   # Start the blockchain network
   ./network.sh up
   
   # Deploy chaincode
   ./deployCC.sh
   ```

4. **Install API dependencies**
   ```bash
   cd api
   npm install
   ```

5. **Start the REST API**
   ```bash
   npm start
   ```

6. **Launch the frontend**
   ```bash
   cd ../frontend
   npm install
   npm start
   ```

## (Theoritical) Usage

### For Talents
1. Access the web interface
2. Submit academic or professional credential requests
3. Monitor approval status
4. View verified credentials

### For Institutions
1. Log in with institutional credentials
2. Review pending credential requests
3. Approve or revoke credentials
4. Manage institutional credential records

### For Companies
1. Access company portal
2. Verify candidate credentials
3. Review credential authenticity
4. Issue professional certifications

## API Documentation

### Credential Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/credentials/academic` | Create academic credential |
| POST | `/credentials/professional` | Create professional credential |
| PUT | `/credentials/{id}/approve` | Approve credential |
| PUT | `/credentials/{id}/revoke` | Revoke credential |
| DELETE | `/credentials/{id}` | Delete credential |
| GET | `/credentials/{id}` | Retrieve credential by ID |
| GET | `/credentials/all` | Retrieve all credentials |
| PUT | `/credentials/{id}/skills` | Update credential skills |
| PUT | `/credentials/{id}/name` | Update credential name |

## Smart Contracts

### Credential Lifecycle Operations

| Function | Description |
|----------|-------------|
| `InitLedger` | Initialize the blockchain ledger |
| `CreateAcademicCredential` | Create new academic credential |
| `CreateProfessionalCredential` | Create new professional credential |
| `UpdateVerificationStatus` | Update credential verification status |
| `GetAllCredentials` | Query all credentials |
| `CredentialExists` | Check if credential exists |
| `DeleteTalentCredential` | Remove credential from ledger |

### Credential Attributes

- **CredentialID**: Unique identifier
- **FirstName/LastName**: Credential holder's name
- **TalentID**: Unique talent identifier
- **Skills**: Array of skills/competencies
- **VerificationStatus**: Current verification state
- **VerifiedBy**: Verifying institution/organization

## Performance Results

### Throughput vs. Latency Analysis
- **Optimal Performance**: Up to 5 concurrent transactions
- **Saturation Point**: Performance degrades with >50 concurrent transactions
- **Transaction Size Impact**: Negligible effect on performance (32B - 512B tested)

### Key Findings
- Network overhead and consensus mechanisms dominate performance
- Payload size has minimal impact on transaction processing
- Typical blockchain throughput vs. latency tradeoff observed

## License

This project is provided for academic purposes. Feel free to explore, modify, or build upon it.

---

**Note**: This system is designed for demonstration and research purposes. For production deployment, additional security measures, scalability optimizations, and comprehensive testing are recommended.
