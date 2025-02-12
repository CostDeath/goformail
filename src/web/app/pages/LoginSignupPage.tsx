"use client"

import Card from "@/components/card";
import LoginForm from "@/components/loginSignup/loginForm";
import {useState} from "react";
import SignUpForm from "@/components/loginSignup/signUpForm";

export default function LoginSignupPage() {
  const [isLogin, setIsLogin] = useState(true);
  const switchCard = () => {
    setIsLogin(!isLogin);
  }
  return (
      <div className="grid grid-rows-auto items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
        <div className="row-start-1">
          {isLogin && (
              <>
                <Card><LoginForm/></Card>
                <div className="grid grid-cols-2 py-5 px-5">
                  <p>Don&apos;t have an account?</p>
                  <button className="text-blue-500 underline font-bold" onClick={switchCard}>Sign up!</button>
                </div>
              </>
          )}

          {!isLogin && (
              <>
                  <Card><SignUpForm/></Card>
                <div className="grid grid-cols-2 py-5 px-5">
                  <p>Already have an account?</p>
                  <button className="text-blue-500 underline font-bold" onClick={switchCard}>Log in!</button>
                </div>
              </>
          )}
        </div>
      </div>
  );
}
