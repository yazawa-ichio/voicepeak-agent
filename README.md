# VOICEPEAK Agent

VOICEPEAKのベータ機能を利用した音声再生用のボットの試作です。  
現状、再生まで三秒ぐらいかかったり、多重起動が出来ないので実用性はそこまでです。  

## 実行

```sh
#ボットの起動
APP_PATH="C:\\Program Files\\VOICEPEAK 6ナレーターセット\\voicepeak.exe"
go run bot/*.go -a "${APP_PATH}" -d "./tmp"
```

```sh
#再生リクエスト
curl -G --data-urlencode text='Hello' --data-urlencode narrator='Female Child' localhost:21952/speak
```


