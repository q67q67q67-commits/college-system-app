import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'providers/auth_provider.dart';
import 'theme/app_theme.dart';
import 'screens/login_screen.dart';
import 'screens/main_shell.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const NCApp());
}

class NCApp extends StatelessWidget {
  const NCApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => AuthProvider()..loadFromStorage(),
      child: MaterialApp(
        title: 'Narxoz College',
        theme: appTheme,
        home: Consumer<AuthProvider>(
          builder: (_, auth, __) {
            if (auth.isLoggedIn) return const MainShell();
            return const LoginScreen();
          },
        ),
      ),
    );
  }
}
