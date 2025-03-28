import Button from "./Button";
import Logo from "./Logo";
import NavButton from "./NavButton";

export default function NavBar() {
  return (
    <nav className="flex justify-between items-center px-6 py-4 bg-white/50 backdrop-blur-2xl shadow-md">
      {/* Left Side: Logo */}
      <Logo />

      {/* Right Side: Buttons */}
      <div className="flex space-x-4">
        <NavButton to="/login" variant="outline">Login</NavButton>
        <NavButton to="/register">Sign Up</NavButton>
      </div>
    </nav>
  );
}