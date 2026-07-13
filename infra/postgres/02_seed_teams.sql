INSERT INTO teams (name, city, abbreviation)
VALUES
    ('Dream', 'Atlanta', 'ATL'),
    ('Sky', 'Chicago', 'CHI'),
    ('Sun', 'Connecticut', 'CON'),
    ('Wings', 'Dallas', 'DAL'),
    ('Valkyries', 'Golden State', 'GSV'),
    ('Fever', 'Indiana', 'IND'),
    ('Aces', 'Las Vegas', 'LVA'),
    ('Sparks', 'Los Angeles', 'LAS'),
    ('Lynx', 'Minnesota', 'MIN'),
    ('Liberty', 'New York', 'NYL'),
    ('Mercury', 'Phoenix', 'PHX'),
    ('Fire', 'Portland', 'POR'),
    ('Storm', 'Seattle', 'SEA'),
    ('Tempo', 'Toronto', 'TOR'),
    ('Mystics', 'Washington', 'WAS')
ON CONFLICT (abbreviation) DO NOTHING;