import Button from "./Button";
import Logo from "./Logo";

export default function NavBar() {
  return (
    <nav className="flex justify-between items-center px-6 py-4 bg-white/50 backdrop-blur-2xl shadow-md">
      {/* Left Side: Logo */}
      <Logo />

      {/* Right Side: Buttons */}
      <div className="flex space-x-4">
        <Button variant="outline">Login</Button>
        <Button>Sign Up</Button>
      </div>
    </nav>
  );
}