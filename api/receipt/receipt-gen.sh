
echo "auto gen receipt.go at receipt"
# -o output file location
# -p expected package name
# yaml file
goapi-gen -o receipt.go -p receipt.gen.go receipt-api.yaml
