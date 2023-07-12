create table cars (
    RegistrationPlate TEXT PRIMARY KEY ,
    Model TEXT NOT NULL ,
    Purpose TEXT NOT NULL CHECK ( Purpose in ('SHARING', 'TAXI', 'DELIVERY') ),
    ManufactureYear INT NOT NULL ,
    Mileage INT NOT NULL CHECK ( Mileage >= 0 )
);
