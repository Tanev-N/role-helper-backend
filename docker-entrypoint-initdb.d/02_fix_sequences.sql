
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));
SELECT setval('characters_id_seq', (SELECT COALESCE(MAX(id), 1) FROM characters));
SELECT setval('session_players_id_seq', (SELECT COALESCE(MAX(id), 1) FROM session_players));

