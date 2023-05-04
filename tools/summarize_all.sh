file_path=$1

for i in {4..6}
do
  bash tools/summarize.sh ./res/no_rec/$i-$i >> $file_path
  echo "" >> $file_path
  bash tools/summarize.sh ./res/rec/$i-$i >> $file_path
  echo -e "\n\n\n" >> $file_path
done

bash tools/summarize.sh ./res/no_rec/8-6 >> $file_path
echo "" >> $file_path
bash tools/summarize.sh ./res/rec/8-6 >> $file_path
