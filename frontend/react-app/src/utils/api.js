import API_BASE_URL from "../config";

export const checkSession = async () => {
    try {
        const response = await fetch(`${API_BASE_URL}/auth/isLoggedIn`, {
            method: "GET",
            credentials: "include",
        });
        if(!response.ok){
            throw new Error("Failed to check session");
        }
        const isLoggedIn = await response.json();
        
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