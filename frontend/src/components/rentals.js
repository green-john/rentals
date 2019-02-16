import { $http } from "./http";
import $auth from "./auth";

export default {
    loadAllApartments(filters) {
        let finalUrl = '/apartments?';
        if (filters.floorAreaMeters) {
            finalUrl += `floorAreaMeters=${filters.floorAreaMeters}&`;
        }

        if (filters.pricePerMonthUSD) {
            finalUrl += `pricePerMonthUSD=${filters.pricePerMonthUSD}&`;
        }

        if (filters.roomCount) {
            finalUrl += `roomCount=${filters.roomCount}`;
        }

        return $http.get(
            finalUrl, {
                headers: {Authorization: $auth.getToken()}
            }
        ).then(response => {
            return response.data;
        }).catch(err => {
            alert(err);
            throw err;
        })
    },

    newApartment(apartmentData) {
        console.log(apartmentData);
        return $http.post('/apartments', apartmentData, {
            headers: {Authorization: $auth.getToken()}
        }).then(response => {
            return response.data;
        }).catch(err => {
            alert(err);
            throw err;
        })
    },

    changeAvailability(apartmentId, newAvailability) {
        return $http.patch('/apartments/' + apartmentId, {
            'available': newAvailability
        }, {
            headers: {Authorization: $auth.getToken()}
        });
    }
}