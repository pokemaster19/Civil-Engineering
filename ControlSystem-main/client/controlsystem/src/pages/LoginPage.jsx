import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import api, { API_BASE_URL } from "../api/axiosInstance";
import styles from "../css/LoginPage.module.css";

const LoginPage = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState(null);
    const navigate = useNavigate();
    const { login } = useAuth();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);

        try {
            const response = await api.axiosInstance.post("/auth/login", {
                email,
                password,
            });

            const { token, user_id, role_id } = response.data;

            login(token, user_id, role_id);
            navigate("/");
        } catch (err) {
            console.error("Ошибка входа:", err);
            setError("Неверный логин или пароль");
        }
    };

    return (
        <div className={styles.background}>
            <div className={styles.loginContainer}>
                <h2>Войдите в свой аккаунт</h2>
                <form onSubmit={handleSubmit}>
                    <div className={styles.formGroup}>
                        <input
                            type="text"
                            id="email"
                            name="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            placeholder="Email"
                        />
                    </div>
                    <div className={styles.formGroup}>
                        <input
                            type="password"
                            id="password"
                            name="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            placeholder="Пароль"
                        />
                    </div>
                    {error && <p className={styles.error}>{error}</p>}
                    <button type="submit" className={styles.loginButton}>
                        Войти
                    </button>
                </form>
            </div>
        </div>
    );
};

export default LoginPage;
