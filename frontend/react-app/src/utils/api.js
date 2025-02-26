import API_BASE_URL from "../config";
import axios from "axios";

export const fetchLoggedinUser = async () => {
    try {
        const response = await axios.get(`${API_BASE_URL}/auth/loggedInUser`, {
            withCredentials: true
        });
        const user = response.data;
        
        return user;
    }catch(error){
        console.error("Error checking session: ", error);

        return null;
    }
}

export const logoutUser = async(csrfToken) => {
    try {
        await axios.post(`${API_BASE_URL}/auth/logout`, {},
            {
                withCredentials: true,
                headers: {
                    'X-CSRF-Token': csrfToken
                }
            }
        );

    }catch(error){
        console.error("Error logging out: ", error);
        throw error;
    }
};