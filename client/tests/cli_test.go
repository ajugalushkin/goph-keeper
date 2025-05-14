package tests

//func TestClientRegisterLoginCli(t *testing.T) {
//	email := gofakeit.Email()
//	password := gofakeit.Password(true, true, true, true, false, 8)
//
//	tests := []struct {
//		name    string
//		args    []string
//		fixture string
//	}{
//		{"start root",
//			[]string{}, "root.golden"},
//		{"start auth",
//			[]string{"auth"}, "auth-no-args.golden"},
//		{"start auth register",
//			[]string{"auth", "register", "--email", email, "--password", password}, "register.golden"},
//		{"start auth login",
//			[]string{"auth", "login", "--email", email, "--password", password}, "login.golden"},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			output, err := runBinary(tt.args)
//
//			if err != nil {
//				t.Fatalf("failed to run binary: %v", err.Error())
//			}
//
//			if *update {
//				writeFixture(t, tt.fixture, output)
//			}
//
//			actual := string(output)
//
//			expected := loadFixture(t, tt.fixture)
//
//			if !reflect.DeepEqual(actual, expected) {
//				t.Fatalf("actual = %s, expected = %s", actual, expected)
//			}
//		})
//	}
//}
//
//func TestClientKeepCreateCli(t *testing.T) {
//	email := gofakeit.Email()
//	password := gofakeit.Password(true, true, true, true, false, 8)
//
//	tests := []struct {
//		name    string
//		args    []string
//		fixture string
//	}{
//		{"start keep",
//			[]string{"keep"}, "keep-no-args.golden"},
//		{"start auth register",
//			[]string{"auth", "register", "--email", email, "--password", password}, "register.golden"},
//		{"start auth login",
//			[]string{"auth", "login", "--email", email, "--password", password}, "login.golden"},
//		//{"start auth register",
//		//	[]string{"auth", "register", "--email", email, "--password", password}, "register.golden"},
//		//{"start auth login",
//		//	[]string{"auth", "login", "--email", email, "--password", password}, "login.golden"},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			output, err := runBinary(tt.args)
//
//			if err != nil {
//				t.Fatalf("failed to run binary: %v", err.Error())
//			}
//
//			if *update {
//				writeFixture(t, tt.fixture, output)
//			}
//
//			actual := string(output)
//
//			expected := loadFixture(t, tt.fixture)
//
//			if !reflect.DeepEqual(actual, expected) {
//				t.Fatalf("actual = %s, expected = %s", actual, expected)
//			}
//		})
//	}
//}

//func TestMain(m *testing.M) {
//	err := os.Chdir("..")
//	if err != nil {
//		fmt.Printf("could not change dir: %v", err)
//		os.Exit(1)
//	}
//
//	dir, err := os.Getwd()
//	if err != nil {
//		fmt.Printf("could not get current dir: %v", err)
//	}
//
//	binaryPath = filepath.Join(dir, binaryName)
//
//	os.Exit(m.Run())
//}
//
//var update = flag.Bool("update", false, "update golden files")
//
//var binaryName = "client"
//
//var binaryPath = ""
//
//func fixturePath(t *testing.T, fixture string) string {
//	_, filename, _, ok := runtime.Caller(0)
//	if !ok {
//		t.Fatalf("problems recovering caller information")
//	}
//
//	return filepath.Join(filepath.Dir(filename), fixture)
//}
//
//func writeFixture(t *testing.T, fixture string, content []byte) {
//	err := os.WriteFile(fixturePath(t, fixture), content, 0644)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func loadFixture(t *testing.T, fixture string) string {
//	content, err := os.ReadFile(fixturePath(t, fixture))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	return string(content)
//}
//
//func runBinary(args []string) ([]byte, error) {
//	cmd := exec.Command(binaryPath, args...)
//	return cmd.CombinedOutput()
//}

//func TestCreateItemStream_Success(t *testing.T) {
//	cfg := config.Client{
//		Address: ":8080",
//		Timeout: time.Hour,
//		Retries: 3,
//	}
//	log := slog.New(
//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//	)
//
//	temp, err := os.Create("test_bin.txt")
//	defer os.Remove("test_bin.txt")
//	require.NoError(t, err)
//
//	stat, err := temp.Stat()
//	require.NoError(t, err)
//
//	fileName := stat.Name()
//
//	_, err = temp.WriteString(gofakeit.Letter())
//	require.NoError(t, err)
//
//	// Register a new user for authentication
//	authClient := auth.NewAuthClient(auth.GetAuthConnection(log, cfg))
//
//	login := gofakeit.Email()
//	password := gofakeit.Password(true, true, true, true, true, 8)
//
//	err = authClient.Register(context.Background(), login, password)
//	require.NoError(t, err)
//
//	// Log in the user
//	newToken, err := authClient.Login(context.Background(), login, password)
//	require.NoError(t, err)
//
//	keepClient := NewKeeperClient(GetKeeperConnection(log, cfg.Address, newToken))
//
//	resp, err := keepClient.CreateItemStream(context.Background(), fileName, fileName)
//	require.NoError(t, err)
//	assert.Equal(t, fileName, resp.GetName())
//}
