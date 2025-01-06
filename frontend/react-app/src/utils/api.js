import API_BASE_URL from "../config";
import axios from "axios";

export const checkSession = async () => {
    try {
        const response = await axios.get(`${API_BASE_URL}/auth/isLoggedIn`, {
            withCredentials: true
        });
        const isLoggedIn = response.data;
        
        return isLoggedIn;
    }catch(error){
        console.error("Error chekcing session: ", error);

        return false;
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