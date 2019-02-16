<template>
  <v-dialog v-model="showModal" max-width="600px">
    <v-card>
      <v-card-title>
        <span class="headline">New Apartment</span>
      </v-card-title>
      <v-card-text>
        <v-container grid-list-md>
          <v-layout wrap>
            <v-flex xs12 sm12 md12>
              <v-text-field
                      label="Name"
                      hint="Name of the apartment"
                      v-model.trim="newApartmentData.name"
                      persistent-hint
                      required>
              </v-text-field>
            </v-flex>
            <v-flex xs12 sm12 md12>
              <v-textarea
                      label="Description"
                      v-model.trim="newApartmentData.description"
              ></v-textarea>
            </v-flex>
            <v-flex xs12 sm12 md4>
              <v-text-field
                      label="Floor Area*"
                      hint="Area in m2"
                      persistent-hint
                      v-model.number="newApartmentData.floorAreaMeters"
                      required
              ></v-text-field>
            </v-flex>
            <v-flex xs12 sm12 md4>
              <v-text-field
                      label="Price per month*"
                      hint="Price in USD"
                      persistent-hint
                      v-model.trim="newApartmentData.pricePerMonthUSD"
                      required
              ></v-text-field>
            </v-flex>
            <v-flex xs12 sm12 md4>
              <v-text-field
                      label="Room count*"
                      hint="Number of rooms"
                      persistent-hint
                      v-model.number="newApartmentData.roomCount"
                      required
              ></v-text-field>
            </v-flex>
            <v-flex xs6 sm6 md6>
              <v-text-field
                      label="Longitude*"
                      v-model.number="newApartmentData.longitude"
                      required
              ></v-text-field>
            </v-flex>
            <v-flex xs6 sm6 md6>
              <v-text-field
                      label="Latitude*"
                      v-model.number="newApartmentData.latitude"
                      required
              ></v-text-field>
            </v-flex>
            <v-flex v-if="isAdmin()" xs12 sm12 md12>
              <v-select
                      :items="allAdminsAndRealtors"
                      label="Realtor*"
                      v-model.number="newApartmentData.realtorId"
                      required
              ></v-select>
            </v-flex>
            <v-flex xs12 sm12 md12>
              <v-checkbox
                      label="Available*"
                      v-model="newApartmentData.available"
                      required
              ></v-checkbox>
            </v-flex>
          </v-layout>
        </v-container>
        <small>*indicates required field</small>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" flat @click="showModal = false">Close</v-btn>
        <v-btn color="blue darken-1" flat @click="createApartment()">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
    import $rentals from './rentals';
    import $users from "./users";

    export default {
        name: "NewApartment",
        props: {
            showModal: Boolean,
        },
        data() {
            return {
                newApartmentData: {
                    name: null,
                    description: null,
                    realtorId: null,
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                    latitude: null,
                    longitude: null,
                    available: null,
                },

                allUsers: [],
            }
        },
        methods: {
            tryLoadingAllUsers() {
                $users.getAllUsers().then(res => {
                    this.allUsers = res;
                }).catch(err => {
                    console.log(err);
                })
            },

            createApartment() {
                // If is not set, set it to myself
                if (!this.newApartmentData.realtorId) {
                    this.newApartmentData.realtorId = this.userData.id;
                }

                $rentals.newApartment(this.newApartmentData).then(res => {
                    alert(`Apartment created with id ${res.id}`);
                    this.clearApartmentData();
                    this.loadApartments();
                    this.showModal = false;
                }).catch(err => {
                    this.newApartmentMessage = `Error: ${err}`;
                });
            },

            clearApartmentData() {
                this.newApartmentData = {
                    name: null,
                    description: null,
                    realtorId: null,
                    floorAreaMeters: null,
                    pricePerMonthUSD: null,
                    roomCount: null,
                    latitude: null,
                    longitude: null,
                    available: null,
                };
                this.newApartmentMessage = "";
            },
        }
    }
</script>

<style scoped>

</style>