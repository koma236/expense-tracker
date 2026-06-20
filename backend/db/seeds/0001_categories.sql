-- 初期カテゴリ。アプリ初回セットアップ時に投入する。
INSERT INTO categories (name, type) VALUES
  -- 支出
  ('食費',       'expense'),
  ('日用品',     'expense'),
  ('交通費',     'expense'),
  ('住居費',     'expense'),
  ('水道光熱費', 'expense'),
  ('娯楽費',     'expense'),
  ('交際費',     'expense'),
  ('医療費',     'expense'),
  ('その他',     'expense'),
  -- 収入
  ('給与',       'income'),
  ('賞与',       'income'),
  ('副収入',     'income'),
  ('その他',     'income');
