# etherium-poc

etherium-pocはブロックチェーンアプリケーションを試しに動かしてみるアプリケーションです

![](./images/image_1.png)

参考書籍: [ブロックチェーン実践入門](https://books.google.com/books?id=-toBEAAAQBAJ&printsec=copyright)

なお書籍の技術は若干古い上、関連技術のバージョンアップに伴う破壊的変更により書籍の通りに実行しても動作しないため、下記のリプレースを行っています

* geth => hardhat
* web3.js => ethers
* html, Javascript => React

## usage

### ブロックチェーンを起動する
```bash
cd deploy-contract
npm install(初回のみ)
npm start
npm run deploy(別のターミナルで実行)
```

### フロントエンドアプリケーションを起動する
```
cd frontend
npm install(初回のみ)
npm start
```

投票ボタンを押すとスマートコントラクトが発行される

Ctrl+Cで起動したブロックチェーンを終了するとリセットされる

![](./images/image_2.png)

### CLIでコントラクトを実行する

フロントエンドでボタンを押した時と同様の挙動をCLIでも再現できる

```
cd cli
go run main.go
```