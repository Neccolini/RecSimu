for i in {4..6}
do
    rm ./res/rec/$i-$i/*
    rm ./res/no_rec/$i-$i/*
done
rm ./res/rec/8-6/*
rm ./res/no_rec/8-6/*

for i in {1..10} 
do
for j in {1..12} 
do
for s in {4..6} 
do
    go run main.go run -i ./examples/rec/$s-$s/$j.jsonc >> res/rec/$s-$s/log${s}_${s}_$j
    go run main.go run -i ./examples/no_rec/$s-$s/$j.jsonc >> res/no_rec/$s-$s/log${s}_${s}_$j
done
    go run main.go run -i ./examples/rec/8-6/$j.jsonc >> res/rec/8-6/log8_6_$j
    go run main.go run -i ./examples/no_rec/8-6/$j.jsonc >> res/no_rec/8-6/log8_6_$j
done
done
