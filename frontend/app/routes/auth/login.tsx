import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import { Link, useNavigate } from "react-router";
import Button from "~/components/ui/Button";
import { useAuth } from "~/lib/auth";
import { api } from "~/lib/http";

export default function Login() {
  const { user, setUser, setAccessToken } = useAuth();
  let navigate = useNavigate();

  const login = useMutation({
    mutationFn: (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      
      const formData = new FormData(e.target as HTMLFormElement);
      const userData = Object.fromEntries(formData.entries()) as {
        username: string;
        password: string;
      };
      return api.post("/users/login", userData);
    },
    onSuccess: (data) => {
      setUser({
        id: data.data.id,
        username: data.data.username,
        profilePicture: data.data.profile_picture,
        createdAt: data.data.created_at,
      });
      setAccessToken(data.data.access);

      navigate("/hub");
    },
  });

  // automatically redirect users if they're already logged in
  useEffect(() => {
    if (user) {
      navigate("/hub");
    }
  }, []);

  return (
    <div className="flex size-full min-h-min flex-col items-center justify-center overflow-y-auto p-12">
      <form className="w-full max-w-md" onSubmit={login.mutate}>
        <h1 className="text-black text-center font-bold text-3xl">login</h1>
        {login.isError ? (
          <h2 className="text-red-500 text-center font-bold text-xl">{login.error.message}</h2>
        ) : null}
        <div className="mt-6">
          <label
            className="block items-center text-sm font-medium leading-5 text-black"
            htmlFor="username"
          >
            username
          </label>
          <input
            className="block w-full mt-1 px-3 py-2 text-black placeholder:text-gray-300 focus:border-rose-400 focus:outline-none focus:ring-4 focus:ring-rose-50 shadow-sm border border-slate-300 rounded-md"
            aria-label="username"
            name="username"
            id="username"
            placeholder="username"
            step="1"
            type="text"
            required
          />
        </div>

        <div className="mt-6">
          <label
            className="block items-center text-sm font-medium leading-5 text-black"
            htmlFor="password"
          >
            password
          </label>
          <input
            className="block w-full mt-1 px-3 py-2 text-black placeholder:text-gray-300 focus:border-rose-400 focus:outline-none focus:ring-4 focus:ring-rose-50 shadow-sm border border-slate-300 rounded-md"
            aria-label="password"
            name="password"
            id="password"
            placeholder="password"
            step="1"
            type="password"
            required
          />
        </div>

        <Button
          className="mt-12 block w-full"
          size="lg"
          type="submit"
          isLoading={login.isPending}
        >
          login
        </Button>

        <div className="text-black text-sm mt-4">
          {"don't have an account? "}
          <Link className="text-rose-400" to="/register">register here</Link>
          .
        </div>
      </form>
    </div>
  );
}