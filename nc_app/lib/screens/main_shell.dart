import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../services/api_service.dart';
import 'login_screen.dart';
import 'home_screen.dart';
import 'schedule_screen.dart';
import 'grades_screen.dart';
import 'groups_screen.dart';
import 'library_screen.dart';
import 'forum_screen.dart';
import 'events_screen.dart';
import 'files_screen.dart';
import 'profile_screen.dart';
import 'director_screen.dart';
import 'regulations_screen.dart';
import 'map_screen.dart';

enum NavItem {
  home('Главная', Icons.home),
  schedule('Расписание', Icons.calendar_today),
  grades('Оценки', Icons.grade),
  groups('Группы', Icons.groups),
  teacherGrades('Выставить оценку', Icons.edit),
  teacherHomework('ДЗ к парам', Icons.assignment),
  library('Библиотека', Icons.menu_book),
  forum('Форум', Icons.forum),
  events('События', Icons.event),
  files('Мои файлы', Icons.folder),
  director('Профиль директора', Icons.person),
  regulations('Регламент', Icons.description),
  map('Карта', Icons.map),
  profile('Профиль', Icons.person_outline);

  final String title;
  final IconData icon;
  const NavItem(this.title, this.icon);
}

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _index = 0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final auth = context.read<AuthProvider>();
      ApiService.onUnauthorized = () async {
        await auth.logout();
        if (context.mounted) Navigator.of(context).pushAndRemoveUntil(MaterialPageRoute(builder: (_) => const LoginScreen()), (_) => false);
      };
    });
  }

  List<NavItem> _navItems(AuthProvider auth) {
    final role = auth.role ?? '';
    final items = <NavItem>[NavItem.home, NavItem.schedule];
    if (role == 'student') items.add(NavItem.grades);
    if (auth.isTeacher) {
      items.add(NavItem.groups);
      items.add(NavItem.teacherGrades);
      items.add(NavItem.teacherHomework);
    }
    items.addAll([NavItem.library, NavItem.forum, NavItem.events, NavItem.files]);
    if (auth.showDirectorProfile) items.add(NavItem.director);
    items.addAll([NavItem.regulations, NavItem.map, NavItem.profile]);
    return items;
  }

  Widget _screen(NavItem item) {
    switch (item) {
      case NavItem.home:
        return const HomeScreen();
      case NavItem.schedule:
        return const ScheduleScreen();
      case NavItem.grades:
        return const GradesScreen();
      case NavItem.groups:
        return const GroupsScreen();
      case NavItem.teacherGrades:
        return const TeacherGradesScreen();
      case NavItem.teacherHomework:
        return const TeacherHomeworkScreen();
      case NavItem.library:
        return const LibraryScreen();
      case NavItem.forum:
        return const ForumScreen();
      case NavItem.events:
        return const EventsScreen();
      case NavItem.files:
        return const FilesScreen();
      case NavItem.director:
        return const DirectorScreen();
      case NavItem.regulations:
        return const RegulationsScreen();
      case NavItem.map:
        return const MapScreen();
      case NavItem.profile:
        return const ProfileScreen();
    }
  }

  @override
  Widget build(BuildContext context) {
    final auth = context.watch<AuthProvider>();
    final items = _navItems(auth);
    if (_index >= items.length) _index = 0;
    final current = items[_index];
    final isWide = MediaQuery.of(context).size.width >= 600;

    return Scaffold(
      appBar: AppBar(
        title: Text(current.title),
        actions: [
          Text('${auth.fullName ?? ""} (${auth.role ?? ""})', style: const TextStyle(fontSize: 12)),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async {
              await auth.logout();
              if (context.mounted) Navigator.of(context).pushAndRemoveUntil(MaterialPageRoute(builder: (_) => const LoginScreen()), (_) => false);
            },
          ),
        ],
      ),
      body: Row(
        children: [
          if (isWide)
            NavigationRail(
              selectedIndex: _index,
              onDestinationSelected: (i) => setState(() => _index = i),
              labelType: NavigationRailLabelType.all,
              destinations: items.map((e) => NavigationRailDestination(icon: Icon(e.icon), label: Text(e.title))).toList(),
            ),
          Expanded(child: _screen(current)),
        ],
      ),
      bottomNavigationBar: isWide ? null : BottomNavigationBar(
        currentIndex: _index,
        onTap: (i) => setState(() => _index = i),
        type: BottomNavigationBarType.fixed,
        items: items.map((e) => BottomNavigationBarItem(icon: Icon(e.icon), label: e.title)).toList(),
      ),
    );
  }
}
