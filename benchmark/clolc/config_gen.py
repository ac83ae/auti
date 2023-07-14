import argparse

import yaml


def generate_fablo_config(
        output_filename,
        chaincode_name,
        chaincode_dir,
        num_orderers,
        orderer_type,
        num_orgs,
        num_auditors,
        num_peers,
        fabric_version,
        desired_fabric_version,
        desired_fabric_ca_version,
):
    """Generate Fablo config."""
    config = {
        "$schema": "https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/schema.json",
        "global": {"fabricVersion": fabric_version, "tls": True},
        "orgs": [],
        "channels": [{"name": "mychannel", "orgs": []}],
        "chaincodes": [
            {
                "name": chaincode_name,
                "version": "0.1.0",
                "lang": "golang",
                "channel": "mychannel",
                "endorsement": "",
                "directory": chaincode_dir,
            }
        ],
        "hooks": {
            "postGenerate": f"perl -i -pe 's/_VERSION={fabric_version}/_VERSION={desired_fabric_version}/g' ./fablo"
                            f"-target/fabric-docker/.env; perl -i -pe 's/FABRIC_CA_VERSION=1.5.0/FA"
                            f"BRIC_CA_VERSION={desired_fabric_ca_version}/g' ./fablo-target/fabric-docker/.env"
        },
    }

    endorsement_list = []

    for i in range(1, num_orgs + 1):
        org_name = f"Org{i}"
        org_domain = f"org{i}.example.com"
        config["orgs"].append(
            {
                "organization": {"name": org_name, "domain": org_domain},
                "peer": {"instances": num_peers},
            }
        )
        peers_list = [f"peer{j}" for j in range(num_peers)]
        config["channels"][0]["orgs"].append(
            {"name": org_name, "peers": peers_list}
        )
        endorsement_list.append(f"'{org_name}MSP.member'")

    for i in range(1, num_auditors + 1):
        aud_name = f"Aud{i}"
        aud_domain = f"aud{i}.example.com"
        config["orgs"].append(
            {
                "organization": {"name": aud_name, "domain": aud_domain},
                "peer": {"instances": num_peers},
            }
        )
        peers_list = [f"peer{j}" for j in range(num_peers)]
        config["channels"][0]["orgs"].append(
            {"name": aud_name, "peers": peers_list}
        )
        endorsement_list.append(f"'{aud_name}MSP.member'")

    config["orgs"].append(
        {
            "organization": {"name": "com", "domain": "com.example.com"},
            "orderers": [
                {
                    "groupName": "group1",
                    "prefix": "orderer",
                    "type": orderer_type,
                    "instances": num_orderers,
                }
            ],
        }
    )

    config["chaincodes"][0]["endorsement"] = (
            "AND(" + ", ".join(endorsement_list) + ")"
    )

    with open(output_filename, "w") as file:
        yaml.dump(config, file)


def main():
    """Main function to handle command line arguments."""
    parser = argparse.ArgumentParser(description="Generate Fablo config")
    parser.add_argument(
        "--output_filename",
        default="fablo_config.yaml",
        type=str,
        help="Output filename",
    )
    parser.add_argument(
        "--chaincode_name",
        default="",
        type=str,
        help="Chaincode name",
    )
    parser.add_argument(
        "--chaincode_dir",
        default="",
        type=str,
        help="Chaincode directory",
    )
    parser.add_argument(
        "--num_orderers",
        default=3,
        type=int,
        help="Number of orderers",
    )
    parser.add_argument(
        "--orderer_type",
        default="raft",
        type=str,
        help="Orderer type: solo or raft, kafka not supported yet",
    )
    parser.add_argument(
        "--num_orgs",
        default=1,
        type=int,
        help="Number of organizations",
    )
    parser.add_argument(
        "--num_peers",
        default=1,
        type=int,
        help="Number of peers",
    )
    parser.add_argument(
        "--num_auditors",
        default=0,
        type=int,
        help="Number of auditors",
    )
    parser.add_argument(
        "--fabric_version",
        default="2.4.3",
        type=str,
        help="Current Fabric version",
    )
    parser.add_argument(
        "--desired_fabric_version",
        default="2.5.3",
        type=str,
        help="Desired Fabric version",
    )
    parser.add_argument(
        "--desired_fabric_ca_version",
        default="1.5.6",
        type=str,
        help="Desired Fabric CA version",
    )

    args = parser.parse_args()

    generate_fablo_config(
        args.output_filename,
        args.chaincode_name,
        args.chaincode_dir,
        args.num_orderers,
        args.orderer_type,
        args.num_orgs,
        args.num_auditors,
        args.num_peers,
        args.fabric_version,
        args.desired_fabric_version,
        args.desired_fabric_ca_version,
    )


if __name__ == "__main__":
    main()
