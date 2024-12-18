## エンティティ（テーブル）一覧

大きく分けて、以下の実体があります。

1. ユーザー関連 (Users)
2. 記事関連 (Articles, Categories, Tags, Article_Tags)
3. コメント (Comments)
4. いいね(Likes)
5. 認証(Users表で対応可能)
6. その他付随情報 (e.g. ログインセッションはトークンベースでアプリケーション側で管理するため、DB上では特にテーブルを用意しないケースもある)

### テーブル一覧

- `users`：ユーザー情報を格納するテーブル
- `categories`：記事分類カテゴリを格納するテーブル
- `tags`：タグを格納するテーブル
- `articles`：記事本体を格納するテーブル
- `article_tags`：記事とタグの多対多関係を保持する中間テーブル
- `comments`：コメントを格納するテーブル
- `likes`：記事に対する「いいね」を格納するテーブル

## テーブル定義詳細

### users テーブル

| カラム名     | 型                | NULL許可 | 説明                            |
|--------------|-------------------|----------|---------------------------------|
| id           | UUID(またはINT AUTO INCREMENT) | NOT NULL | 主キー                           |
| email        | VARCHAR(255)      | NOT NULL | メールアドレス (ユニーク制約)    |
| username     | VARCHAR(50)       | NOT NULL | ユーザー名 (ユニーク制約)        |
| password_hash| VARCHAR(255)      | NOT NULL | ハッシュ化したパスワード         |
| created_at   | TIMESTAMP (UTC)   | NOT NULL | ユーザー作成日時 (DEFAULT CURRENT_TIMESTAMP) |

**インデックス・制約**:
- `UNIQUE(email)`
- `UNIQUE(username)`
- 必要に応じて検索性能向上のため `INDEX(email)`

### categories テーブル

| カラム名   | 型               | NULL許可 | 説明                        |
|------------|------------------|----------|-----------------------------|
| id         | UUID or INT AUTO | NOT NULL | 主キー                      |
| name       | VARCHAR(100)     | NOT NULL | カテゴリ名(ユニーク制約)     |

**インデックス・制約**:
- `UNIQUE(name)`

### tags テーブル

| カラム名   | 型               | NULL許可 | 説明                        |
|------------|------------------|----------|-----------------------------|
| id         | UUID or INT AUTO | NOT NULL | 主キー                      |
| name       | VARCHAR(100)     | NOT NULL | タグ名(ユニーク制約)         |

**インデックス・制約**:
- `UNIQUE(name)`

### articles テーブル

| カラム名     | 型               | NULL許可 | 説明                            |
|--------------|------------------|----------|---------------------------------|
| id           | UUID or INT AUTO | NOT NULL | 主キー                           |
| title        | VARCHAR(255)     | NOT NULL | 記事タイトル                     |
| content      | TEXT             | NOT NULL | 記事コンテンツ本体               |
| author_id    | UUID or INT      | NOT NULL | 記事作成者(users.idへのFK)       |
| category_id  | UUID or INT      | NOT NULL | カテゴリID(categories.idへのFK)  |
| created_at   | TIMESTAMP        | NOT NULL | 作成日時(Default CURRENT_TIMESTAMP) |
| updated_at   | TIMESTAMP        | NOT NULL | 更新日時(Default CURRENT_TIMESTAMP ON UPDATE) |

**インデックス・制約**:
- `FOREIGN KEY (author_id) REFERENCES users(id)`
- `FOREIGN KEY (category_id) REFERENCES categories(id)`
- 検索性能向上のためtitle, category_idなどにインデックスを考慮
- full-text searchを行う場合はcontentやtitleに対してFull-Text IndexをサポートするDBで検討

※ いいね数(like_count)はarticlesテーブルにキャッシュするか、likesテーブルから集計するかは設計方針による。もし集計コスト削減のためにキャッシュを持つなら、`like_count` INT DEFAULT 0をarticlesテーブルに持ち、likesテーブルの更新でトリガーないしアプリ側でインクリメントする運用があり得る。

### article_tags テーブル (多対多中間テーブル)

| カラム名     | 型          | NULL許可 | 説明                     |
|--------------|-------------|----------|--------------------------|
| article_id   | UUID or INT | NOT NULL | 記事ID(articles.idへのFK)|
| tag_id       | UUID or INT | NOT NULL | タグID(tags.idへのFK)    |

**インデックス・制約**:
- `PRIMARY KEY (article_id, tag_id)`
- `FOREIGN KEY (article_id) REFERENCES articles(id)`
- `FOREIGN KEY (tag_id) REFERENCES tags(id)`
- `INDEX(article_id)`
- `INDEX(tag_id)`

### comments テーブル

| カラム名    | 型               | NULL許可 | 説明                          |
|-------------|------------------|----------|-------------------------------|
| id          | UUID or INT AUTO | NOT NULL | 主キー                         |
| article_id  | UUID or INT      | NOT NULL | 紐付く記事ID(articles.idへのFK) |
| author_id   | UUID or INT      | NOT NULL | コメント投稿者(users.idへのFK)   |
| content     | TEXT             | NOT NULL | コメント本文                   |
| created_at  | TIMESTAMP        | NOT NULL | 作成日時(Default CURRENT_TIMESTAMP)|

**インデックス・制約**:
- `FOREIGN KEY (article_id) REFERENCES articles(id)`
- `FOREIGN KEY (author_id) REFERENCES users(id)`
- `INDEX(article_id)`
- コメント一覧取得時にarticle_idで検索することが多いので`INDEX(article_id)`は有用

### likes テーブル

| カラム名    | 型               | NULL許可 | 説明                           |
|-------------|------------------|----------|--------------------------------|
| id          | UUID or INT AUTO | NOT NULL | 主キー                          |
| article_id  | UUID or INT      | NOT NULL | いいね対象記事(articles.idへのFK)|
| user_id     | UUID or INT      | NOT NULL | いいねしたユーザー(users.idへのFK)|
| created_at  | TIMESTAMP        | NOT NULL | 作成日時(Default CURRENT_TIMESTAMP)|

**インデックス・制約**:
- `FOREIGN KEY (article_id) REFERENCES articles(id)`
- `FOREIGN KEY (user_id) REFERENCES users(id)`
- `UNIQUE(article_id, user_id)`：同一ユーザーが同一記事に重複して「いいね」できないようにする
- `INDEX(article_id), INDEX(user_id)`

## 補足

- 認証はJWTを想定しているため、特定のセッションテーブルは用意していないが、必要に応じてrefresh tokenや有効セッションを管理するテーブルを設けることもある。
- RSS用のエンドポイントは記事データを元に生成するだけなので、特に追加テーブルは不要。
- 検索用パラメータ（カテゴリー、タグ、タイトル・内容全文検索）に対して、インデックスやFull-Text Search対応が必要な場合はDB製品依存で対応。
- `id`カラムはユニーク性が必要だが、DBごとに異なる識別子戦略を使用できる。PostgreSQLならUUID型を使ったり、`SERIAL`や`BIGSERIAL`を使ったりできる。MySQLなら`AUTO_INCREMENT`を使用、あるいはUUIDを生成して格納することも可能。
- 更新系(APIでpatchやdelete)のアクセスは、セキュリティ制約(bearerAuth)で特定ユーザーに紐づくレコードのみ操作が可能なアプリケーションロジックが別途必要。
