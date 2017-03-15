create keyspace eas with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table eas.ingested_session(id UUID, PRIMARY KEY(id));
