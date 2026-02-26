package com.huatai.release;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@MapperScan("com.huatai.release.mapper")
public class ReleasePlatformApplication {

    public static void main(String[] args) {
        SpringApplication.run(ReleasePlatformApplication.class, args);
    }
}
