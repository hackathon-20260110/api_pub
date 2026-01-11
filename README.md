# api
## run
実行方法を示す。次のコマンドを実行することで8080ポートでgoのプロセスが立ち上がる。
また、同時に5432ポートでpostgresqlのプロセスが立ち上がる。
```bash
docker compose up --build -d
```

## development
開発時にはdevcontainerを使用することを推奨する。
VScodeまたはcursor(などのvscodeに準じたエディタ)でdevcontainerを開くことで、コンテナ内での開発が可能となる。

### devcontainer内での開発
以下の内容はdevcontainer内での開発方法である。
Makefileをタスクランナーとして利用している。

次のコマンドで利用可能なコマンド一覧を表示することができる。
```bash
make help # または　make
```

開発ビルドは `make dev`で行う。
`air`というgoのホットリロードツールを使用しているためファイルの変更が即時反映される。

### db
データベースにはPostgreSQLを使用している。
ローカル開発環境においてはdocker composeで立ち上げている。
リモートのDBでは **Neon**というサービスを利用予定である。
