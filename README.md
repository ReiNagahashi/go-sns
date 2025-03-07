

# GO_SNS  

## 概要  
GO_SNSは、Goを使用して実装した軽量なSNSアプリケーションです。このプロジェクトの目的は、設計、サーバー構築、データベース、セキュリティ、オブジェクト指向プログラミングなど、これまで学んだソフトウェア開発の知識を実践的にアウトプットすることにあります。そのため、フレームワークは使用せず、純粋にGo言語の機能のみで実装しています。  


## プロジェクトの目的  
- ソフトウェアの設計から運用までを体系的に実践。  
- MVCアーキテクチャを通じて、オブジェクト指向プログラミング（OOP）の理解を深める。  
- CSRF対策や暗号処理の実装を通じたセキュリティの復習。  
- Dockerを使用し、環境依存の問題を解決して高速にアプリケーションをデプロイ可能にする。  
- テスト駆動開発（TDD）と自動デプロイを実現するためのCI/CD環境を構築。  
- 並列処理や非同期処理の実装手法を習得。  

## 主要なポイント  
1. **クラス設計**: アプリケーション層、データアクセス層、データ層を明確化するため、UML図を使用した設計を実施。![gosns_uml](https://github.com/user-attachments/assets/cbcf1e84-cd59-4948-9bb3-77f3b916b804)
2. **データベースマイグレーション**: データベーススキーマを迅速かつ直感的に更新可能なマイグレーションスクリプトを実装。  
3. **データシーディング**: API機能の迅速な動作確認を可能にするデータシーディングを導入。  
4. **単体テストとリファクタリング**: 各モジュールの依存性を分離したコードへリファクタリング。  
5. **セッションベース認証**: ユーザーの認証にセッションを活用。  
6. **CI/CD**: 一部の単体テストをリモートプッシュ時に自動実行するCI/CDパイプラインを構築。  


## 機能要件  
### 最低要件 (Minimum Viable Product: MVP)  
- ユーザー登録: メールアドレスとハッシュ化されたパスワードでアカウントを作成。  
- セッション管理: ログインしたユーザーに一意のセッションIDを付与。  
- 投稿機能: 投稿の作成・削除が可能。
- ランダムな投稿データのシーディング：並列処理により、非同期的に新たな投稿を作成・表示。

### （実装/拡張予定） 
- 「いいね」機能: 他ユーザーの投稿に「いいね」を付与。  
- リアルタイムチャット機能：ユーザー同士がコネクションレスにメッセージのやり取りが可能なチャット機能の実装
- 通知機能: 「いいね」やメッセージ受信時に通知を送信。 
- 水平スケーリングの検討(複数のサーバが追加された分散サーバサイドシステム)
   1. 10~100万件の投稿が毎分スケジュールされる場合のデータの分散方法
   2. インフルエンサーのような人気ユーザーの投稿のようなデータによるデータベース負荷の軽減をするためのキャッシュの設計
   3. 常に新しいサーバが追加および削除される場合、この水平スケーリングを行う際のデプロイメントの管理方法


## アプリの動作手順  

1. **Dockerのインストール**  
   Dockerをインストールし、Dockerアカウントにログインします。  
   [公式サイト](https://docs.docker.com/desktop/)  

2. **リポジトリのクローン**  
   以下のコマンドでリポジトリをクローンします:  
   ```bash
   git clone https://github.com/ReiNagahashi/go-sns.git
   cd go-sns
   ```

3. **アプリのビルド & 実行**  
   以下のコマンドでアプリをビルドし、実行します:  
   ```bash
   docker-compose up --build
   ```  
   ※ 初回の実行ではサーバーが立ち上がるまで約30秒かかります。  

