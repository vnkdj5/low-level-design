Design In memory DB including indexes, select query support etc.Design In memory DB including indexes, select query support etc.




## Problem Statement

You are tasked with building an **in-memory database** that mimics the behavior of a lightweight relational database. The database should support the following functionalities:

1. **Table Management**:  
   - Create tables with specified schemas (column names and data types).
   - Support dynamic addition of indexes on table columns to optimize search queries.

2. **CRUD Operations**:  
   - Insert records into a table with unique keys.
   - Retrieve records by their keys.
   - Query data using simple conditions or complex conditions with logical operators like `AND` and `OR`.  
   - Delete records based on their keys.

3. **Querying Support**:  
   - Retrieve specific attributes of records using a `SELECT` query.  
   - Support filtering with conditions such as `=`, `!=`, `<`, `>`, `<=`, `>=`, and `BETWEEN`.  
   - Combine multiple conditions using `AND` or `OR`.

4. **Indexing**:  
   - Create indexes on table columns to optimize queries involving those columns.

### Example Use Case

A developer needs to manage a simple **"Users" database** in memory for a high-speed, low-latency application. The database should allow the following operations:

- Create a table `users` with columns: `id`, `name`, `age`, and `city`.
- Insert user records into the table.
- Retrieve specific user information (e.g., names of users above the age of 25 living in "Mumbai").
- Delete a user record based on their `id`.
- Create an index on the `city` column to optimize queries involving city-based filtering.


## Future scope
1. Add capacity for the number od records
2. Add TTL Support
3. Improve quering
4. Accept query as string and parse it internally