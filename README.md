# 名前
 
論文自動配信 Slack Bot
 
登録したキーワードがタイトルに含まれる論文を，自動配信してくれる Slack Bot
 
 
# 特徴
 
- 論文検索に使いたいキーワードを__複数__登録できます．
- 登録したキーワードがタイトルに含まれる論文(arxiv で検索）を自動配信してくれます．

# 開発理由

読みたい論文を調べなくても自動配信してくれるアプリがあればよいなと思い，作りました．

# 開発期間

2022/4 ~ 

# デモ
 
- 論文検索に使いたいキーワードを登録
https://user-images.githubusercontent.com/63027348/164585306-ee882d0c-5354-4121-85e6-760387cf6b9a.mp4

- 論文検索に使いたいキーワードから削除
https://user-images.githubusercontent.com/63027348/164586179-d0a37c81-1694-4b6c-bc0f-7a73488836d9.mp4

- 論文検索に使いたいキーワード一覧
https://user-images.githubusercontent.com/63027348/164585522-64028952-2ddb-4c20-8e36-711bdad51795.mp4

- 論文の自動配信
https://user-images.githubusercontent.com/63027348/164585831-30c606a6-e3a9-4b53-8c7d-595abd7d0493.mp4



# 使用技術

Go, Slack API in Go, arxiv API, PostgreSQL, Heroku
 
<!--
# Requirement
 
"hoge"を動かすのに必要なライブラリなどを列挙する
 
* huga 3.5.2
* hogehuga 1.0.2
 
# Installation
 
Requirementで列挙したライブラリなどのインストール方法を説明する
 
```bash
pip install huga_package
```
 
# Usage
 
DEMOの実行方法など、"hoge"の基本的な使い方を説明する
 
```bash
git clone https://github.com/hoge/~
cd examples
python demo.py
```
 
# Note
 
注意点などがあれば書く
 
# Author
 
作成情報を列挙する
-->