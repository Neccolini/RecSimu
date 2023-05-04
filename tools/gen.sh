for s in {4..6}
do
for i in {1..12}
do
  cp examples/rec/${s}-${s}/example${s}_${s}.jsonc examples/rec/$s-$s/$i.jsonc
  cp examples/rec/${s}-${s}/example${s}_${s}.jsonc examples/no_rec/$s-$s/$i.jsonc
  python3 ./tools/rec.py ./examples/rec/$s-$s/$i.jsonc $i
  python3 ./tools/rec2.py ./examples/no_rec/$s-$s/$i.jsonc $i
done
done

cp examples/rec/8-6/* examples/no_rec/8-6/

for i in {1..12}
do
  python3 ./tools/rec.py ./examples/rec/8-6/$i.jsonc $i
  python3 ./tools/rec2.py ./examples/no_rec/8-6/$i.jsonc $i
done