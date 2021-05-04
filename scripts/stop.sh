# shellcheck disable=SC2006

id1=`pgrep auth`
id2=`pgrep sub`
id3=`pgrep main`

kill $id1
kill $id2
kill $id3