#/bin/bash
#TODO check for existings gobs first
for f in *; do
    if [ -d "$f" ]; then
        cd $f
        go mod graph | modv $1 ../results
        cd ..
    fi
done
cat results/*.gob > results.gob
