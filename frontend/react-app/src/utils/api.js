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
        console.error("Error chekcing session: ", error);

        return null;
    }
}

export const logoutUser = async() => {
    try {
        const response = await fetch(`${API_BASE_URL}/auth/logout`, {
            method: "POST",
            credentials: "include",
        });
        if(!response.ok){
            throw new Error("Failed to logout");
        }
    }catch(error){
        console.error("Error logging out: ", error);
        throw error;
    }
};