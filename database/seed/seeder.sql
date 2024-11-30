INSERT INTO Users (UserUsername, UserEmail, UserPassword) VALUES
('user1', 'user1@example.com', 'password1'),
('user2', 'user2@example.com', 'password2'),
('user3', 'user3@example.com', 'password3'),
('user4', 'user4@example.com', 'password4'),
('user5', 'user5@example.com', 'password5'),
('user6', 'user6@example.com', 'password6'),
('user7', 'user7@example.com', 'password7'),
('user8', 'user8@example.com', 'password8'),
('user9', 'user9@example.com', 'password9'),
('user10', 'user10@example.com', 'password10');


INSERT INTO Posts (PostTitle, PostBody, LikeCount, created_at, updated_at, Users_UserID) VALUES
('Post Title 1', 'This is the body of post 1.', 5, datetime('now'), datetime('now'), 1),
('Post Title 2', 'This is the body of post 2.', 3, datetime('now'), datetime('now'), 2),
('Post Title 3', 'This is the body of post 3.', 8, datetime('now'), datetime('now'), 3),
('Post Title 4', 'This is the body of post 4.', 2, datetime('now'), datetime('now'), 4),
('Post Title 5', 'This is the body of post 5.', 0, datetime('now'), datetime('now'), 5),
('Post Title 6', 'This is the body of post 6.', 4, datetime('now'), datetime('now'), 6),
('Post Title 7', 'This is the body of post 7.', 1, datetime('now'), datetime('now'), 7),
('Post Title 8', 'This is the body of post 8.', 7, datetime('now'), datetime('now'), 8),
('Post Title 9', 'This is the body of post 9.', 9, datetime('now'), datetime('now'), 9),
('Post Title 10', 'This is the body of post 10.', 6, datetime('now'), datetime('now'), 10);
