solc contracts/token/ERC721.sol --bin --abi -o ../abi/ --overwrite

abigen --sol contracts/token/ERC721.sol --out ./go-contracts/token/erc721.go --pkg token --type erc721 --bin ../abi/AccountDid.bin