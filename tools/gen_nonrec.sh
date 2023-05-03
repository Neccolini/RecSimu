for i in 10 20 30 40 50;
do
for j in {1..9}
do
  for k in {0..5}
  do
    go run main.go gen -c 50000 -d -f ./test.jsonc -t "random $i" -d -r 0.0$j
    echo "$i 0.0$j" >> test_res
    go run main.go run -i ./test.jsonc >> test_res
  done
done
done
