import 'package:flutter/material.dart';
import 'core/theme.dart';
import 'features/auth/login_screen.dart';
import 'features/shell/main_shell.dart';
import 'services/auth_state.dart';

class NcApp extends StatefulWidget {
  const NcApp({super.key});

  @override
  State<NcApp> createState() => _NcAppState();
}

class _NcAppState extends State<NcApp> {
  final AuthState _auth = AuthState();

  @override
  void initState() {
    super.initState();
    _auth.loadFromStorage();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _auth,
      builder: (context, _) {
        return MaterialApp(
          title: 'Narxoz College',
          theme: buildNcTheme(),
          debugShowCheckedModeBanner: false,
          home: _auth.isAuthenticated
              ? MainShell(authState: _auth)
              : LoginScreen(authState: _auth),
        );
      },
    );
  }
}
