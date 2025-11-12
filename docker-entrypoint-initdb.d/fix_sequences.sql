
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));
SELECT setval('characters_id_seq', (SELECT COALESCE(MAX(id), 1) FROM characters));
SELECT setval('game_players_id_seq', (SELECT COALESCE(MAX(id), 1) FROM game_players));

