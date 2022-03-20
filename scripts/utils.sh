function ensure_ok() {
    for ((i=0;i<10;i++)) ;
    do
        echo "$*"
        eval "$@"
        if [ $? -eq 0 ]; then
            break
        fi
        echo "$*"
    done
}