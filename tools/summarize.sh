dir_path=$1

for file in "$dir_path"/*
do
  echo "$file"
  python3 tools/rec_sum.py $file
  echo ""
done