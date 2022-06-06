# Rosetta images for Docker

**Elrond Network runs on a sharded architecture** - transaction, data and network sharding are leveraged. 

In the Rosetta implementation, we've decided to provide a single-shard perspective to the API consumer. That is, **one Rosetta instance** would observe **a single _regular_ shard** of the network (plus the _metachain_) - the shard is selected by the owner of the instance.

The Rosetta deployment for Elrond takes the shape of two Docker images (Elrond Rosetta and Elrond Observer), plus a Docker Compose definition to orchestrate the `1 + 1 + 1 + 1 = 4` containers: 

 - one Elrond Rosetta instance in **online mode**
 - one Elrond Rosetta instance in **offline mode**
 - one Elrond observer for a chosen regular shard
 - one Elrond observer for the _metachain_ (necessary for some pieces of information such as [ESDT](https://docs.elrond.com/developers/esdt-tokens) properties)
 
This `1 + 1 + 1 + 1 = 4` setup is usually referred to as an **Elrond Rosetta Squad**.

## Prerequisites

### Clone this repository

```
cd $HOME
git clone https://github.com/ElrondNetwork/rosetta.git
```

### Build the images

Below, we build all the images (including for _testnet_ and _devnet_).

```
cd $HOME/rosetta

docker image build . -t elrond-rosetta:latest -f ./docker/Rosetta.dockerfile

docker image build . -t elrond-rosetta-observer-testnet:latest -f ./docker/ObserverTestnet.dockerfile
docker image build . -t elrond-rosetta-observer-devnet:latest -f ./docker/ObserverDevnet.dockerfile
docker image build . -t elrond-rosetta-observer-mainnet:latest -f ./docker/ObserverMainnet.dockerfile
```

### Prepare folders on host

The following script prepares the required folder structure on host:

```
cd $HOME/rosetta

./prepare_host.sh ${HOME}/rosetta
```

### Generate keys for observers

The following script generates the node keys, required by the observers (chosen shard, plus metachain):

```
cd $HOME/rosetta

./generate_keys.sh ${HOME}/rosetta/keys
```

## Run rosetta

### Run on testnet

```
cd $HOME/rosetta

export ROSETTA_IMAGE=elrond-rosetta:latest
export OBSERVER_IMAGE=elrond-rosetta-observer-testnet:latest
export DATA_FOLDER=${HOME}/rosetta/testnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

### Run on devnet

```
cd $HOME/rosetta

export ROSETTA_IMAGE=elrond-rosetta:latest
export OBSERVER_IMAGE=elrond-rosetta-observer-devnet:latest
export DATA_FOLDER=${HOME}/rosetta/devnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

### Run on mainnet

```
cd $HOME/rosetta

export ROSETTA_IMAGE=elrond-rosetta:latest
export OBSERVER_IMAGE=elrond-rosetta-observer-mainnet:latest
export DATA_FOLDER=${HOME}/rosetta/mainnet
export KEYS_FOLDER=${HOME}/rosetta/keys

docker compose --file ./docker-compose.yml up --detach
```

## Update rosetta

Update the repository (repositories):

```
cd $HOME/rosetta
git pull origin
```

Stop the running containers:

```
docker stop rosetta-observer-0-1
docker stop rosetta-observer-metachain-1
docker stop rosetta-rosetta-1
docker stop rosetta-rosetta-offline-1
```

Re-build the images as described above, then run the containers again.
