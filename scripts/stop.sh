# shellcheck disable=SC2006

id1=`pgrep auth_service`
id2=`pgrep sub_service`
id3=`pgrep chat_service`
id5=`pgrep kudago_service`
id4=`pgrep main_service`

kill $id1
kill $id2
kill $id3
kill $id5
kill $id4

rm chat.out
rm sub.out
rm kudago.out
rm auth.out
rm main.out