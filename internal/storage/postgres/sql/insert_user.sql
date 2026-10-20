INSERT INTO users(id,email,password)
VALUES(gen_random_uuid(),$1,$2)
RETURNING id;