INSERT INTO org.company (name, display_name) VALUES
  ('acme-analytics',  'Acme Analytics Inc.'),
  ('blue-horizon',    'Blue Horizon Research'),
  ('vertex-data',    'Vertex Data Solutions');

INSERT INTO auth.user (username, password, email, role) VALUES
  ('jvalentine', crypt('password123', gen_salt('bf')), 'jill.valentine@acme-analytics.com',   'analist'),
  ('credfield',  crypt('password123', gen_salt('bf')), 'chris.redfield@acme-analytics.com',   'analist'),
  ('rchambers',  crypt('password123', gen_salt('bf')), 'rebecca.chambers@blue-horizon.com',   'analist'),
  ('bburton',    crypt('password123', gen_salt('bf')), 'barry.burton@blue-horizon.com',        'analist'),
  ('awestker',   crypt('password123', gen_salt('bf')), 'albert.wesker@vertex-data.com',        'manager');

INSERT INTO org.analist (user_id, first_name, last_name, company_id) VALUES
  (1, 'Jill',     'Valentine', 1),
  (2, 'Chris',    'Redfield',  1),
  (3, 'Rebecca',  'Chambers',  2),
  (4, 'Barry',    'Burton',    2),
  (5, 'Albert',   'Wesker',    3);

INSERT INTO ngin.storage (name, uri) VALUES
  ('Q1 Financial Reports',        's3://acme-data/finance/q1-2026/reports.db'),
  ('Customer Churn Dataset',      's3://acme-data/ml/churn/v3/churn.db'),
  ('Market Research PDFs',        'gs://blue-horizon-docs/market-research/2026/research.db'),
  ('ESG Compliance Documents',    'gs://blue-horizon-docs/esg/compliance/compliance.db'),
  ('Sales Forecast Models',       's3://vertex-prod/forecasting/sales-2026/forecast.db'),
  ('Contract Archive',            'azblob://vertex-legal/contracts/active/contracts.db');

INSERT INTO ngin.storage_analist (storage_id, analist_id) VALUES
  (1, 1), (1, 2),
  (2, 1),
  (3, 3), (3, 4),
  (4, 3),
  (5, 5),
  (6, 5);

INSERT INTO ngin.model (name, provider, api_key, company_id) VALUES
  ('gpt-4.1',           'openai',    'sk-openai-acme-finance-dev', 1),
  ('claude-3-5-haiku',  'anthropic', 'sk-ant-blue-horizon-ops-dev', 2),
  ('grok-3-mini',       'xai',       'sk-xai-vertex-forecast-dev', 3);

INSERT INTO auth.apikey (hashed_key, user_id, name, active, expiry)
SELECT encode(hmac(raw_key, 'dev-hmac-secret', 'sha256'), 'hex'), user_id, name, active, expiry
FROM (VALUES
  ('m8k_7f3d9a2b-1c4e-4f8a-9b0d-2e5f6a7c8d9e', 1, 'Jill Valentine — primary key',   TRUE,  NOW() + INTERVAL '1 year'),
  ('m8k_2a1b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d', 2, 'Chris Redfield — primary key',   TRUE,  NOW() + INTERVAL '1 year'),
  ('m8k_9e8d7c6b-5a4f-3e2d-1c0b-9a8f7e6d5c4b', 3, 'Rebecca Chambers — primary key', TRUE,  NOW() + INTERVAL '6 months'),
  ('m8k_4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e', 5, 'Albert Wesker — primary key',    TRUE,  NOW() + INTERVAL '1 year'),
  ('m8k_1a2b3c4d-5e6f-7a8b-9c0d-ef1234567890', 1, 'Legacy CI key',                  FALSE, NOW() - INTERVAL '1 day')
) AS t(raw_key, user_id, name, active, expiry);

INSERT INTO ngin.notebook (analist_id, title, chat) VALUES
  (1, 'Jill Notebook 01', '[
    {"role": "user",      "content": "Sample prompt 1 for analist 1."},
    {"role": "assistant", "content": "Sample response 1 for analist 1."}
  ]'),
  (1, 'Jill Notebook 02', '[
    {"role": "user",      "content": "Sample prompt 2 for analist 1."},
    {"role": "assistant", "content": "Sample response 2 for analist 1."}
  ]'),
  (1, 'Jill Notebook 03', '[
    {"role": "user",      "content": "Sample prompt 3 for analist 1."},
    {"role": "assistant", "content": "Sample response 3 for analist 1."}
  ]'),
  (1, 'Jill Notebook 04', '[
    {"role": "user",      "content": "Sample prompt 4 for analist 1."},
    {"role": "assistant", "content": "Sample response 4 for analist 1."}
  ]'),
  (1, 'Jill Notebook 05', '[
    {"role": "user",      "content": "Sample prompt 5 for analist 1."},
    {"role": "assistant", "content": "Sample response 5 for analist 1."}
  ]'),
  (1, 'Jill Notebook 06', '[
    {"role": "user",      "content": "Sample prompt 6 for analist 1."},
    {"role": "assistant", "content": "Sample response 6 for analist 1."}
  ]'),
  (1, 'Jill Notebook 07', '[
    {"role": "user",      "content": "Sample prompt 7 for analist 1."},
    {"role": "assistant", "content": "Sample response 7 for analist 1."}
  ]'),
  (1, 'Jill Notebook 08', '[
    {"role": "user",      "content": "Sample prompt 8 for analist 1."},
    {"role": "assistant", "content": "Sample response 8 for analist 1."}
  ]'),
  (1, 'Jill Notebook 09', '[
    {"role": "user",      "content": "Sample prompt 9 for analist 1."},
    {"role": "assistant", "content": "Sample response 9 for analist 1."}
  ]'),
  (1, 'Jill Notebook 10', '[
    {"role": "user",      "content": "Sample prompt 10 for analist 1."},
    {"role": "assistant", "content": "Sample response 10 for analist 1."}
  ]'),
  (1, 'Jill Notebook 11', '[
    {"role": "user",      "content": "Sample prompt 11 for analist 1."},
    {"role": "assistant", "content": "Sample response 11 for analist 1."}
  ]'),
  (1, 'Jill Notebook 12', '[
    {"role": "user",      "content": "Sample prompt 12 for analist 1."},
    {"role": "assistant", "content": "Sample response 12 for analist 1."}
  ]'),
  (1, 'Jill Notebook 13', '[
    {"role": "user",      "content": "Sample prompt 13 for analist 1."},
    {"role": "assistant", "content": "Sample response 13 for analist 1."}
  ]'),
  (1, 'Jill Notebook 14', '[
    {"role": "user",      "content": "Sample prompt 14 for analist 1."},
    {"role": "assistant", "content": "Sample response 14 for analist 1."}
  ]'),
  (1, 'Jill Notebook 15', '[
    {"role": "user",      "content": "Sample prompt 15 for analist 1."},
    {"role": "assistant", "content": "Sample response 15 for analist 1."}
  ]'),
  (1, 'Jill Notebook 16', '[
    {"role": "user",      "content": "Sample prompt 16 for analist 1."},
    {"role": "assistant", "content": "Sample response 16 for analist 1."}
  ]'),
  (1, 'Jill Notebook 17', '[
    {"role": "user",      "content": "Sample prompt 17 for analist 1."},
    {"role": "assistant", "content": "Sample response 17 for analist 1."}
  ]'),
  (1, 'Jill Notebook 18', '[
    {"role": "user",      "content": "Sample prompt 18 for analist 1."},
    {"role": "assistant", "content": "Sample response 18 for analist 1."}
  ]'),
  (1, 'Jill Notebook 19', '[
    {"role": "user",      "content": "Sample prompt 19 for analist 1."},
    {"role": "assistant", "content": "Sample response 19 for analist 1."}
  ]'),
  (1, 'Jill Notebook 20', '[
    {"role": "user",      "content": "Sample prompt 20 for analist 1."},
    {"role": "assistant", "content": "Sample response 20 for analist 1."}
  ]'),
  (1, 'Jill Notebook 21', '[
    {"role": "user",      "content": "Sample prompt 21 for analist 1."},
    {"role": "assistant", "content": "Sample response 21 for analist 1."}
  ]'),
  (1, 'Jill Notebook 22', '[
    {"role": "user",      "content": "Sample prompt 22 for analist 1."},
    {"role": "assistant", "content": "Sample response 22 for analist 1."}
  ]'),
  (1, 'Jill Notebook 23', '[
    {"role": "user",      "content": "Sample prompt 23 for analist 1."},
    {"role": "assistant", "content": "Sample response 23 for analist 1."}
  ]'),
  (2, 'Chris Notebook 01', '[
    {"role": "user",      "content": "Sample prompt 1 for analist 2."},
    {"role": "assistant", "content": "Sample response 1 for analist 2."}
  ]'),
  (2, 'Chris Notebook 02', '[
    {"role": "user",      "content": "Sample prompt 2 for analist 2."},
    {"role": "assistant", "content": "Sample response 2 for analist 2."}
  ]'),
  (2, 'Chris Notebook 03', '[
    {"role": "user",      "content": "Sample prompt 3 for analist 2."},
    {"role": "assistant", "content": "Sample response 3 for analist 2."}
  ]'),
  (2, 'Chris Notebook 04', '[
    {"role": "user",      "content": "Sample prompt 4 for analist 2."},
    {"role": "assistant", "content": "Sample response 4 for analist 2."}
  ]'),
  (2, 'Chris Notebook 05', '[
    {"role": "user",      "content": "Sample prompt 5 for analist 2."},
    {"role": "assistant", "content": "Sample response 5 for analist 2."}
  ]'),
  (2, 'Chris Notebook 06', '[
    {"role": "user",      "content": "Sample prompt 6 for analist 2."},
    {"role": "assistant", "content": "Sample response 6 for analist 2."}
  ]'),
  (2, 'Chris Notebook 07', '[
    {"role": "user",      "content": "Sample prompt 7 for analist 2."},
    {"role": "assistant", "content": "Sample response 7 for analist 2."}
  ]'),
  (2, 'Chris Notebook 08', '[
    {"role": "user",      "content": "Sample prompt 8 for analist 2."},
    {"role": "assistant", "content": "Sample response 8 for analist 2."}
  ]'),
  (2, 'Chris Notebook 09', '[
    {"role": "user",      "content": "Sample prompt 9 for analist 2."},
    {"role": "assistant", "content": "Sample response 9 for analist 2."}
  ]'),
  (2, 'Chris Notebook 10', '[
    {"role": "user",      "content": "Sample prompt 10 for analist 2."},
    {"role": "assistant", "content": "Sample response 10 for analist 2."}
  ]'),
  (2, 'Chris Notebook 11', '[
    {"role": "user",      "content": "Sample prompt 11 for analist 2."},
    {"role": "assistant", "content": "Sample response 11 for analist 2."}
  ]'),
  (2, 'Chris Notebook 12', '[
    {"role": "user",      "content": "Sample prompt 12 for analist 2."},
    {"role": "assistant", "content": "Sample response 12 for analist 2."}
  ]'),
  (2, 'Chris Notebook 13', '[
    {"role": "user",      "content": "Sample prompt 13 for analist 2."},
    {"role": "assistant", "content": "Sample response 13 for analist 2."}
  ]'),
  (2, 'Chris Notebook 14', '[
    {"role": "user",      "content": "Sample prompt 14 for analist 2."},
    {"role": "assistant", "content": "Sample response 14 for analist 2."}
  ]'),
  (2, 'Chris Notebook 15', '[
    {"role": "user",      "content": "Sample prompt 15 for analist 2."},
    {"role": "assistant", "content": "Sample response 15 for analist 2."}
  ]'),
  (2, 'Chris Notebook 16', '[
    {"role": "user",      "content": "Sample prompt 16 for analist 2."},
    {"role": "assistant", "content": "Sample response 16 for analist 2."}
  ]'),
  (2, 'Chris Notebook 17', '[
    {"role": "user",      "content": "Sample prompt 17 for analist 2."},
    {"role": "assistant", "content": "Sample response 17 for analist 2."}
  ]'),
  (2, 'Chris Notebook 18', '[
    {"role": "user",      "content": "Sample prompt 18 for analist 2."},
    {"role": "assistant", "content": "Sample response 18 for analist 2."}
  ]'),
  (2, 'Chris Notebook 19', '[
    {"role": "user",      "content": "Sample prompt 19 for analist 2."},
    {"role": "assistant", "content": "Sample response 19 for analist 2."}
  ]'),
  (2, 'Chris Notebook 20', '[
    {"role": "user",      "content": "Sample prompt 20 for analist 2."},
    {"role": "assistant", "content": "Sample response 20 for analist 2."}
  ]'),
  (2, 'Chris Notebook 21', '[
    {"role": "user",      "content": "Sample prompt 21 for analist 2."},
    {"role": "assistant", "content": "Sample response 21 for analist 2."}
  ]'),
  (2, 'Chris Notebook 22', '[
    {"role": "user",      "content": "Sample prompt 22 for analist 2."},
    {"role": "assistant", "content": "Sample response 22 for analist 2."}
  ]'),
  (2, 'Chris Notebook 23', '[
    {"role": "user",      "content": "Sample prompt 23 for analist 2."},
    {"role": "assistant", "content": "Sample response 23 for analist 2."}
  ]'),
  (3, 'Rebecca Notebook 01', '[
    {"role": "user",      "content": "Sample prompt 1 for analist 3."},
    {"role": "assistant", "content": "Sample response 1 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 02', '[
    {"role": "user",      "content": "Sample prompt 2 for analist 3."},
    {"role": "assistant", "content": "Sample response 2 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 03', '[
    {"role": "user",      "content": "Sample prompt 3 for analist 3."},
    {"role": "assistant", "content": "Sample response 3 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 04', '[
    {"role": "user",      "content": "Sample prompt 4 for analist 3."},
    {"role": "assistant", "content": "Sample response 4 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 05', '[
    {"role": "user",      "content": "Sample prompt 5 for analist 3."},
    {"role": "assistant", "content": "Sample response 5 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 06', '[
    {"role": "user",      "content": "Sample prompt 6 for analist 3."},
    {"role": "assistant", "content": "Sample response 6 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 07', '[
    {"role": "user",      "content": "Sample prompt 7 for analist 3."},
    {"role": "assistant", "content": "Sample response 7 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 08', '[
    {"role": "user",      "content": "Sample prompt 8 for analist 3."},
    {"role": "assistant", "content": "Sample response 8 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 09', '[
    {"role": "user",      "content": "Sample prompt 9 for analist 3."},
    {"role": "assistant", "content": "Sample response 9 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 10', '[
    {"role": "user",      "content": "Sample prompt 10 for analist 3."},
    {"role": "assistant", "content": "Sample response 10 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 11', '[
    {"role": "user",      "content": "Sample prompt 11 for analist 3."},
    {"role": "assistant", "content": "Sample response 11 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 12', '[
    {"role": "user",      "content": "Sample prompt 12 for analist 3."},
    {"role": "assistant", "content": "Sample response 12 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 13', '[
    {"role": "user",      "content": "Sample prompt 13 for analist 3."},
    {"role": "assistant", "content": "Sample response 13 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 14', '[
    {"role": "user",      "content": "Sample prompt 14 for analist 3."},
    {"role": "assistant", "content": "Sample response 14 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 15', '[
    {"role": "user",      "content": "Sample prompt 15 for analist 3."},
    {"role": "assistant", "content": "Sample response 15 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 16', '[
    {"role": "user",      "content": "Sample prompt 16 for analist 3."},
    {"role": "assistant", "content": "Sample response 16 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 17', '[
    {"role": "user",      "content": "Sample prompt 17 for analist 3."},
    {"role": "assistant", "content": "Sample response 17 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 18', '[
    {"role": "user",      "content": "Sample prompt 18 for analist 3."},
    {"role": "assistant", "content": "Sample response 18 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 19', '[
    {"role": "user",      "content": "Sample prompt 19 for analist 3."},
    {"role": "assistant", "content": "Sample response 19 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 20', '[
    {"role": "user",      "content": "Sample prompt 20 for analist 3."},
    {"role": "assistant", "content": "Sample response 20 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 21', '[
    {"role": "user",      "content": "Sample prompt 21 for analist 3."},
    {"role": "assistant", "content": "Sample response 21 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 22', '[
    {"role": "user",      "content": "Sample prompt 22 for analist 3."},
    {"role": "assistant", "content": "Sample response 22 for analist 3."}
  ]'),
  (3, 'Rebecca Notebook 23', '[
    {"role": "user",      "content": "Sample prompt 23 for analist 3."},
    {"role": "assistant", "content": "Sample response 23 for analist 3."}
  ]'),
  (4, 'Barry Notebook 01', '[
    {"role": "user",      "content": "Sample prompt 1 for analist 4."},
    {"role": "assistant", "content": "Sample response 1 for analist 4."}
  ]'),
  (4, 'Barry Notebook 02', '[
    {"role": "user",      "content": "Sample prompt 2 for analist 4."},
    {"role": "assistant", "content": "Sample response 2 for analist 4."}
  ]'),
  (4, 'Barry Notebook 03', '[
    {"role": "user",      "content": "Sample prompt 3 for analist 4."},
    {"role": "assistant", "content": "Sample response 3 for analist 4."}
  ]'),
  (4, 'Barry Notebook 04', '[
    {"role": "user",      "content": "Sample prompt 4 for analist 4."},
    {"role": "assistant", "content": "Sample response 4 for analist 4."}
  ]'),
  (4, 'Barry Notebook 05', '[
    {"role": "user",      "content": "Sample prompt 5 for analist 4."},
    {"role": "assistant", "content": "Sample response 5 for analist 4."}
  ]'),
  (4, 'Barry Notebook 06', '[
    {"role": "user",      "content": "Sample prompt 6 for analist 4."},
    {"role": "assistant", "content": "Sample response 6 for analist 4."}
  ]'),
  (4, 'Barry Notebook 07', '[
    {"role": "user",      "content": "Sample prompt 7 for analist 4."},
    {"role": "assistant", "content": "Sample response 7 for analist 4."}
  ]'),
  (4, 'Barry Notebook 08', '[
    {"role": "user",      "content": "Sample prompt 8 for analist 4."},
    {"role": "assistant", "content": "Sample response 8 for analist 4."}
  ]'),
  (4, 'Barry Notebook 09', '[
    {"role": "user",      "content": "Sample prompt 9 for analist 4."},
    {"role": "assistant", "content": "Sample response 9 for analist 4."}
  ]'),
  (4, 'Barry Notebook 10', '[
    {"role": "user",      "content": "Sample prompt 10 for analist 4."},
    {"role": "assistant", "content": "Sample response 10 for analist 4."}
  ]'),
  (4, 'Barry Notebook 11', '[
    {"role": "user",      "content": "Sample prompt 11 for analist 4."},
    {"role": "assistant", "content": "Sample response 11 for analist 4."}
  ]'),
  (4, 'Barry Notebook 12', '[
    {"role": "user",      "content": "Sample prompt 12 for analist 4."},
    {"role": "assistant", "content": "Sample response 12 for analist 4."}
  ]'),
  (4, 'Barry Notebook 13', '[
    {"role": "user",      "content": "Sample prompt 13 for analist 4."},
    {"role": "assistant", "content": "Sample response 13 for analist 4."}
  ]'),
  (4, 'Barry Notebook 14', '[
    {"role": "user",      "content": "Sample prompt 14 for analist 4."},
    {"role": "assistant", "content": "Sample response 14 for analist 4."}
  ]'),
  (4, 'Barry Notebook 15', '[
    {"role": "user",      "content": "Sample prompt 15 for analist 4."},
    {"role": "assistant", "content": "Sample response 15 for analist 4."}
  ]'),
  (4, 'Barry Notebook 16', '[
    {"role": "user",      "content": "Sample prompt 16 for analist 4."},
    {"role": "assistant", "content": "Sample response 16 for analist 4."}
  ]'),
  (4, 'Barry Notebook 17', '[
    {"role": "user",      "content": "Sample prompt 17 for analist 4."},
    {"role": "assistant", "content": "Sample response 17 for analist 4."}
  ]'),
  (4, 'Barry Notebook 18', '[
    {"role": "user",      "content": "Sample prompt 18 for analist 4."},
    {"role": "assistant", "content": "Sample response 18 for analist 4."}
  ]'),
  (4, 'Barry Notebook 19', '[
    {"role": "user",      "content": "Sample prompt 19 for analist 4."},
    {"role": "assistant", "content": "Sample response 19 for analist 4."}
  ]'),
  (4, 'Barry Notebook 20', '[
    {"role": "user",      "content": "Sample prompt 20 for analist 4."},
    {"role": "assistant", "content": "Sample response 20 for analist 4."}
  ]'),
  (4, 'Barry Notebook 21', '[
    {"role": "user",      "content": "Sample prompt 21 for analist 4."},
    {"role": "assistant", "content": "Sample response 21 for analist 4."}
  ]'),
  (4, 'Barry Notebook 22', '[
    {"role": "user",      "content": "Sample prompt 22 for analist 4."},
    {"role": "assistant", "content": "Sample response 22 for analist 4."}
  ]'),
  (4, 'Barry Notebook 23', '[
    {"role": "user",      "content": "Sample prompt 23 for analist 4."},
    {"role": "assistant", "content": "Sample response 23 for analist 4."}
  ]');
