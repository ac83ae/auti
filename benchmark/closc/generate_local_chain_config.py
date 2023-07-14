import sys

import yaml

LOCAL_CHAIN_CONFIG_FILENAME = "local-chain-config-sc.yaml"


def generate_config(num_chaincodes):
    config = {
        "$schema": "https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/schema.json",
        "global": {"fabricVersion": "2.4.3", "tls": True},
        "hooks": {
            "postGenerate": "perl -i -pe 's/_VERSION=2.4.3/_VERSION=2.5.3/g' ./fablo-target/fabric-docker/.env; perl -i -pe 's/FABRIC_CA_VERSION=1.5.0/FABRIC_CA_VERSION=1.5.6/g' ./fablo-target/fabric-docker/.env"
        },
        "orgs": [
            {
                "organization": {"name": "Org1", "domain": "org1.example.com"},
                "peer": {"instances": 1},
                "orderers": [
                    {
                        "groupName": "group1",
                        "prefix": "orderer",
                        "type": "raft",
                        "instances": 3,
                    }
                ],
            },
            {
                "organization": {"name": "Aud1", "domain": "aud1.example.com"},
                "peer": {"instances": 1},
            },
        ],
        "channels": [
            {
                "name": "mychannel",
                "orgs": [
                    {"name": "Org1", "peers": ["peer0"]},
                    {"name": "Aud1", "peers": ["peer0"]},
                ],
            }
        ],
        "chaincodes": [],
    }

    for i in range(num_chaincodes):
        chaincode = {
            "name": f"auti-local-chain{i}",
            "version": "0.1.0",
            "lang": "golang",
            "channel": "mychannel",
            "endorsement": "AND('Org1MSP.member', 'Aud1MSP.member')",
            "directory": "contract/local_chain",
        }
        config["chaincodes"].append(chaincode)

    return config


def main():
    if len(sys.argv) < 2:
        print("Usage: python generate_config.py <number_of_chaincodes>")
        sys.exit(1)

    num_chaincodes = int(sys.argv[1])
    config = generate_config(num_chaincodes)

    with open(LOCAL_CHAIN_CONFIG_FILENAME, "w") as outfile:
        yaml.dump(config, outfile, default_flow_style=False, sort_keys=False)


if __name__ == "__main__":
    main()
