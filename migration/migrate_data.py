import uuid
import random
import string
from datetime import datetime, timedelta

def random_string(length=10):
    return ''.join(random.choice(string.ascii_lowercase) for _ in range(length))

def random_email():
    return random_string(5) + '@example.com'

def random_datetime():
    now = datetime.now()
    delta = random.randint(0, 365)
    dt = now - timedelta(days=delta)
    return dt.strftime('%Y-%m-%d %H:%M:%S')

NUM_USERS = 100
NUM_CATEGORIES = 10
NUM_TAGS = 30
NUM_ARTICLES = 500
NUM_COMMENTS = 2000
NUM_LIKES = 3000

# UUIDリストを作成してidとして使用
users_ids = [str(uuid.uuid4()) for _ in range(NUM_USERS)]
categories_ids = [str(uuid.uuid4()) for _ in range(NUM_CATEGORIES)]
tags_ids = [str(uuid.uuid4()) for _ in range(NUM_TAGS)]
articles_ids = [str(uuid.uuid4()) for _ in range(NUM_ARTICLES)]
comments_ids = [str(uuid.uuid4()) for _ in range(NUM_COMMENTS)]
likes_ids = [str(uuid.uuid4()) for _ in range(NUM_LIKES)]

file_name = '9_migrate_data.sql'

# ファイルに書き込む
with open(file_name, 'w') as f:
    # 1. Users
    users = []
    for i in range(NUM_USERS):
        user_id = users_ids[i]
        email = f"user{i+1}.{users_ids[i]}@example.com"
        username = f"user{i+1}.{users_ids[i]}"
        password_hash = "hash_" + random_string(20)
        created_at = random_datetime()
        users.append((user_id, email, username, password_hash, created_at))

    f.write("-- +goose Up\n")
    f.write("-- INSERT INTO users\n")
    for u in users:
        f.write(f"INSERT INTO `users` (id, email, username, password_hash, created_at) VALUES ('{u[0]}', '{u[1]}', '{u[2]}', '{u[3]}', '{u[4]}');\n")

    # 2. Categories
    categories = []
    for i in range(NUM_CATEGORIES):
        cat_id = categories_ids[i]
        name = "category_" + random_string(5)
        categories.append((cat_id, name))

    f.write("\n-- INSERT INTO categories\n")
    for c in categories:
        f.write(f"INSERT INTO `categories` (id, name) VALUES ('{c[0]}', '{c[1]}');\n")

    # 3. Tags
    tags = []
    for i in range(NUM_TAGS):
        t_id = tags_ids[i]
        name = "tag_" + random_string(5)
        tags.append((t_id, name))

    f.write("\n-- INSERT INTO tags\n")
    for t in tags:
        f.write(f"INSERT INTO `tags` (id, name) VALUES ('{t[0]}', '{t[1]}');\n")

    # 4. Articles
    articles = []
    for i in range(NUM_ARTICLES):
        a_id = articles_ids[i]
        title = "title_" + random_string(10)
        content = '\n- '.join(random_string(5) for _ in range(10))
        author_id = random.choice(users_ids)
        category_id = random.choice(categories_ids)
        created_at = random_datetime()
        updated_at = created_at
        articles.append((a_id, title, content, author_id, category_id, created_at, updated_at))

    f.write("\n-- INSERT INTO articles\n")
    for a in articles:
        f.write(f"INSERT INTO `articles` (id, title, content, author_id, category_id, created_at, updated_at) VALUES ('{a[0]}', '{a[1]}', '{a[2]}', '{a[3]}', '{a[4]}', '{a[5]}', '{a[6]}');\n")

    # 5. article_tags
    f.write("\n-- INSERT INTO article_tags\n")
    for a in articles:
        tag_count = random.randint(0, 5)
        if tag_count > 0:
            chosen_tags = random.sample(tags_ids, tag_count)
            for ct in chosen_tags:
                f.write(f"INSERT INTO `article_tags` (article_id, tag_id) VALUES ('{a[0]}', '{ct}');\n")

    # 6. Comments
    comments = []
    for i in range(NUM_COMMENTS):
        c_id = comments_ids[i]
        article_id = random.choice(articles_ids)
        author_id = random.choice(users_ids)
        content = "comment_" + random_string(30)
        created_at = random_datetime()
        comments.append((c_id, article_id, author_id, content, created_at))

    f.write("\n-- INSERT INTO comments\n")
    for c in comments:
        f.write(f"INSERT INTO `comments` (id, article_id, author_id, content, created_at) VALUES ('{c[0]}', '{c[1]}', '{c[2]}', '{c[3]}', '{c[4]}');\n")

    # 7. Likes
    like_set = set()
    unique_likes = []
    # 重複避けるためにsetで管理
    while len(unique_likes) < NUM_LIKES:
        article_id = random.choice(articles_ids)
        user_id = random.choice(users_ids)
        if (article_id, user_id) not in like_set:
            like_set.add((article_id, user_id))
            unique_likes.append((article_id, user_id))

    likes = []
    for i, (aid, uid) in enumerate(unique_likes):
        l_id = likes_ids[i]
        created_at = random_datetime()
        likes.append((l_id, aid, uid, created_at))

    f.write("\n-- INSERT INTO likes\n")
    for l in likes:
        f.write(f"INSERT INTO `likes` (id, article_id, user_id, created_at) VALUES ('{l[0]}', '{l[1]}', '{l[2]}', '{l[3]}');\n")
