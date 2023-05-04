print_res() {
echo $1
file_path="./examples/no_rec/$1.jsonc"
res_file_path="./res/no_rec/$1"
for j in {1..20}
do
  go run main.go gen -f $file_path -c 20000 -t "random $1"  -r 0.01 > /dev/null
  go run main.go run -i $file_path -d > res/no_rec/$1

  python3 regex.py ./res/no_rec/$1 $(($1-1))
done
echo ""
}

for i in {9..15}
do
print_res $i
done

for i in {24..48}
do
print_res $i
done
