CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT,
  google_id TEXT UNIQUE,
  avatar_url TEXT,
  bio TEXT,
  website TEXT,
  role TEXT NOT NULL DEFAULT 'user',
  status TEXT NOT NULL DEFAULT 'active',
  follower_count INT DEFAULT 0,
  following_count INT DEFAULT 0,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE sections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  is_nsfw BOOLEAN DEFAULT false,
  is_regional BOOLEAN DEFAULT false,
  icon TEXT,
  post_count INT DEFAULT 0,
  moderator_id UUID REFERENCES users(id) ON DELETE SET NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE challenges (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  theme TEXT NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  status TEXT DEFAULT 'open',
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_challenges_status ON challenges(status, starts_at);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  section_id UUID NOT NULL REFERENCES sections(id) ON DELETE SET NULL,
  title TEXT NOT NULL,
  media_type TEXT NOT NULL,
  media_url TEXT,
  thumb_url TEXT,
  media_key TEXT,
  width INT,
  height INT,
  status TEXT NOT NULL DEFAULT 'processing',
  is_nsfw BOOLEAN DEFAULT false,
  upvotes INT DEFAULT 0,
  downvotes INT DEFAULT 0,
  comment_count INT DEFAULT 0,
  hot_score DOUBLE PRECISION DEFAULT 0,
  challenge_id UUID REFERENCES challenges(id) ON DELETE SET NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_posts_status_time ON posts(status, inserted_at);
CREATE INDEX idx_posts_section_hot ON posts(section_id, hot_score);
CREATE INDEX idx_posts_hot ON posts(hot_score);

CREATE TABLE post_tags (
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag TEXT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (post_id, tag)
);
CREATE INDEX idx_post_tags_tag ON post_tags(tag);

CREATE TABLE votes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  value INT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (user_id, post_id)
);
CREATE INDEX idx_votes_post ON votes(post_id);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
  body TEXT NOT NULL,
  media_url TEXT,
  media_key TEXT,
  upvotes INT DEFAULT 0,
  downvotes INT DEFAULT 0,
  status TEXT DEFAULT 'visible',
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_comments_post_parent ON comments(post_id, parent_id);

CREATE TABLE comment_votes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  value INT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (user_id, comment_id)
);

CREATE TABLE follows (
  follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_type TEXT NOT NULL,
  target_id UUID NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (follower_id, target_type, target_id)
);

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  ref_id UUID,
  actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
  message TEXT,
  read_at TIMESTAMPTZ,
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_notifications_user ON notifications(user_id, read_at);

CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_type TEXT NOT NULL,
  target_id UUID NOT NULL,
  reason TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'pending',
  handled_by UUID REFERENCES users(id) ON DELETE SET NULL,
  handled_at TIMESTAMPTZ,
  note TEXT,
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mod_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  moderator_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  target_type TEXT NOT NULL,
  target_id UUID NOT NULL,
  action TEXT NOT NULL,
  note TEXT,
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE meme_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  image_url TEXT NOT NULL,
  width INT,
  height INT,
  text_boxes JSONB DEFAULT '{}',
  is_official BOOLEAN DEFAULT false,
  status TEXT DEFAULT 'pending',
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE duels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  post_a_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  post_b_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  status TEXT DEFAULT 'open',
  winner_post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_duels_status ON duels(status, ends_at);

CREATE TABLE duel_votes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  duel_id UUID NOT NULL REFERENCES duels(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (user_id, duel_id)
);

CREATE TABLE challenge_entries (
  challenge_id UUID NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (challenge_id, post_id)
);

CREATE TABLE tips (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
  amount_coins INT NOT NULL,
  status TEXT DEFAULT 'completed',
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_coins (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  balance_coins INT DEFAULT 0,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE coin_packages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  coins INT NOT NULL,
  price_idr INT NOT NULL,
  play_sku TEXT NOT NULL UNIQUE,
  is_active BOOLEAN DEFAULT true,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE coin_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  amount_coins INT NOT NULL,
  balance_after INT NOT NULL,
  ref_id UUID,
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_coin_tx_user ON coin_transactions(user_id, inserted_at);

CREATE TABLE post_reactions (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  reaction TEXT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (user_id, post_id)
);

CREATE TABLE bookmarks (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (user_id, post_id)
);

CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  target_type TEXT,
  target_id UUID,
  metadata JSONB DEFAULT '{}',
  inserted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE app_settings (
  key TEXT PRIMARY KEY,
  value JSONB,
  updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
  inserted_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
