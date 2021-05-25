# geneth
Mimetic Learning project

## Basic Idea
In simple terms, Proof of Work (PoW) looks for a nonce that if hashed with the block header returns a value smaller than a threshold.

PoW should not be solvable through a genetic algorithm because hash functions disrupt the locality principle. However, other NP problems do not present such a limitation. For this reason, the original challenge has been modified in maximizing a fitness funtion (not sure if it is NP, but it can be easily replaced).

The PoW difficulty (modifiable in the genesis.json file) now represents the minimum acceptable fitness. The nonce represents the individual. The hash fuction is now a fitness function (more details are present in the form of comments directly in the source code).

## Edited Files
cd ./go-ethereum-master/consensus/ethash

genetic.go -> a simple implementation of  a genetic algorithm

sealer.go, function: mine -> a GA is instantiated in order to find a solution to the current challenge (it is computed from the hahs of the previous block, but it should be updated to use the hash of the current one)

consensus.go -> CalcDifficulty has been modified in order to return always the same difficulty. verifySeal has been modified in order to check that the fitness of the nonce is above the difficulty

The other files in this folder are the original ones

The other files in the other folders are only needed to setup a blockchain network

## Setup
cd ./go-ethereum-master
docker build -t geneth .

# Running 
(cd .. if in the go-ethereum-master folder)

sudo ./cleanup
./setup

//This starts 4 geneth nodes
docker compose up

//Wait for the DAG generation (I did not clean the code from the original PoW setup)

//This deploys a simple smart contract (from node 0)
docker compose -f deployer.yml up

//This increments the smart contarct once.  (from node 0)
docker compose -f tester.yml up

//This reads the current value of the smart contract (from node 1)
docker compose -f loader.yml up

//When finished 
docker compose down
docker compose -f deployer.yml down
docker compose -f tester.yml down
docker compose -f loader.yml down
sudo ./cleanup
