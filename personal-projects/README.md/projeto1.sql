CREATE TABLE trabalhadores (
		trabalhador_id int,
		fist_name VARCHAR(50),
		last_name VARCHAR(50),
		pagamento_hora DECIMAL(5, 2),
		data_ingresso DATE

);

ALTER TABLE trabalhadores
RENAME COLUMN data_ingresso TO email;


ALTER TABLE trabalhadores
MODIFY COLUMN email varchar(50)
AFTER last_name;

INSERT INTO trabalhadores
VALUES(1, "Odilon", "Neves", "odiloneves@gmail.com", 37.15);

INSERT INTO trabalhadores
VALUES
(2, "Amanda", "Neves","mandsrn@gmail.com", 15.00),
(3, "Ricardo", "Bonna","ricbonna@gmail.com", 5.00);

ALTER TABLE trabalhadores
RENAME COLUMN fist_name TO first_name;


SELECT first_name, email
FROM trabalhadores;

SELECT * FROM trabalhadores
WHERE trabalhador_id < 3;

SELECT * FROM trabalhadores;

