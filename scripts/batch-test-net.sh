for ((i=0;i<2;i++))
do
#    for period in 15
#    do
#        sed -i '' "s/PERIOD=[0-9]\{1,\}/PERIOD=$period/g" scripts/test-net.sh
#        echo "======+++++    start $period $i    +++++========="
#        bash scripts/test-net.sh
#        echo "======+++++    end   $period $i    +++++========="
#    done

    for period in 200 400 800
    do
        sed -i '' "s/PERIOD=[0-9]\{1,\}/PERIOD=$period/g" scripts/test-net.sh
        echo "======+++++    start $period $i    +++++========="
        bash scripts/test-net.sh
        echo "======+++++    end   $period $i    +++++========="
    done
done