# Rentals

Manage rental properties

## Spec

*Write an application that manages apartment rentals using a REST API*

* Users must be able to create an account and log in.

* Implement `client`, `realtor` and `admin` role:
   * Clients: browse rentable apartments in a list and on a map.
   * Realtors: client + CRUD all apartments and set the apartment state to available/rented.
   * Admins: realtor +  CRUD realtors, and clients.
   
* Apartments have:
    * Name.
    * Description.
    * Floor area size.
    * Price per month.
    * Number of rooms.
    * Valid geolocation coordinates (either lat/log or geocode).
    * Date added.
    * Associated realtor.

* Apartments are searchable by:
    * Floor area size.
    * Price per month.
    * Number of rooms.
 
- Single-page application. All actions need to be done client side using AJAX,
refreshing the page is not acceptable. Functional UI/UX design is needed. You are
not required to create a unique design, however, do follow best practices to make
the project as functional as possible.

- Bonus: unit and e2e tests.

## Attack plan:

- ~~User creation.~~
- ~~Authentication for user creation.~~
- ~~Apartment creation.~~
- ~~Add authorization to user and apartment creation.~~
- ~~Add read/update/delete tournaments.~~
- ~~Add read/update/delete users.~~
- ~~Search by floor area size, price, rooms.~~
- ~~Write frontend.~~
- ~~Create endpoint to return user info.~~
- Frontend: Only show create to admins/realtors.
- Change available state.
- Validate apartment info.
- Date added for apartments.
- Create client account endpoint.
- Make username unique in database.
- Bonus: create users with admin.
- Bonus: Do geocoding in the frontend.
