for ((i=0;i<2;i++))
do
#    for period in 12 13 14 15 25 50 100 200 400 800 1600
#    do
#        sed -i '' "s/PERIOD=[0-9]\{1,\}/PERIOD=$period/g" scripts/test-all.sh
#        echo "======+++++    start $period $i    +++++========="
#        bash scripts/test-all.sh
#        echo "======+++++    end   $period $i    +++++========="
#    done

    for period in 200 400 800 1600
    do
        sed -i '' "s/PERIOD=[0-9]\{1,\}/PERIOD=$period/g" scripts/test-all-small.sh
        echo "======+++++    start $period $i    +++++========="
        bash scripts/test-all-small.sh
        echo "======+++++    end   $period $i    +++++========="
    done
done