# mall项目构建指南 - 微服务DDD版

本文档提供了如何构建和运行mall电商系统的详细指南。mall是一套电商系统，包括前台商城系统及后台管理系统，基于微服务架构和领域驱动设计(DDD)原则实现，采用Spring Cloud技术栈，并通过Docker容器化部署。

## 目录

- [项目简介](#项目简介)
- [微服务架构与DDD](#微服务架构与DDD)
- [技术架构](#技术架构)
- [环境准备](#环境准备)
- [项目构建步骤](#项目构建步骤)
- [运行各个微服务](#运行各个微服务)
- [Docker部署](#docker部署)
- [常见问题](#常见问题)

## 项目简介

mall项目是一套电商系统，包括前台商城系统及后台管理系统。前台商城系统包含首页门户、商品推荐、商品搜索、商品展示、购物车、订单流程、会员中心、客户服务、帮助中心等模块。后台管理系统包含商品管理、订单管理、会员管理、促销管理、运营管理、内容管理、统计报表、财务管理、权限管理、设置等模块。

新版mall基于微服务架构和DDD原则进行了重构，将系统拆分为多个独立部署的微服务，并按照领域边界进行了清晰划分。

## 微服务架构与DDD

### 领域驱动设计(DDD)简介

领域驱动设计(Domain-Driven Design, DDD)是一种通过将实现与持续进化的模型相连接的方式来满足复杂需求的软件开发方法。mall项目采用DDD的战略设计和战术设计思想，将复杂的业务领域拆分为多个限界上下文(Bounded Context)，并在每个微服务内部应用六边形架构(Hexagonal Architecture)。

### 战略设计

mall系统的核心领域被划分为以下限界上下文，每个限界上下文对应一个或多个微服务：

1. **用户领域(User Domain)**：负责用户注册、认证、授权、会员管理等功能
2. **商品领域(Product Domain)**：负责商品管理、分类、品牌、属性、库存等功能
3. **订单领域(Order Domain)**：负责订单创建、支付、退款、物流等功能
4. **营销领域(Marketing Domain)**：负责促销活动、优惠券、积分、秒杀等功能
5. **内容领域(Content Domain)**：负责CMS、广告、专题、帮助中心等功能
6. **搜索领域(Search Domain)**：负责商品搜索、智能推荐等功能
7. **基础设施(Infrastructure Domain)**：提供通用技术服务支持

### 战术设计

每个微服务内部采用六边形架构(也称为端口与适配器架构)，将领域逻辑与外部依赖解耦：

- **领域层(Domain Layer)**：包含实体(Entity)、值对象(Value Object)、领域服务(Domain Service)、聚合(Aggregate)和聚合根(Aggregate Root)
- **应用层(Application Layer)**：包含应用服务(Application Service)，协调领域对象完成用户用例
- **适配器层(Adapter Layer)**：包含入站适配器(如API控制器)和出站适配器(如数据库仓库实现)
- **基础设施层(Infrastructure Layer)**：提供技术实现，如持久化、消息、缓存等

### 微服务项目组织结构

```
mall-microservice
├── mall-common -- 通用代码和工具类
│   ├── mall-common-core -- 核心通用代码
│   ├── mall-common-web -- Web层通用代码
│   └── mall-common-mybatis -- MyBatis通用代码
├── mall-infrastructure -- 基础设施服务
│   ├── mall-gateway -- API网关服务
│   ├── mall-auth -- 认证授权服务
│   ├── mall-registry -- 服务注册中心
│   └── mall-config -- 配置中心
├── mall-user -- 用户领域微服务
│   ├── mall-user-api -- 用户服务API接口和DTO定义
│   ├── mall-user-application -- 用户应用服务层
│   ├── mall-user-domain -- 用户领域层
│   ├── mall-user-infrastructure -- 用户基础设施层
│   └── mall-user-adapter -- 用户适配器层
├── mall-product -- 商品领域微服务
│   ├── mall-product-api
│   ├── mall-product-application
│   ├── mall-product-domain
│   ├── mall-product-infrastructure
│   └── mall-product-adapter
├── mall-order -- 订单领域微服务
│   ├── mall-order-api
│   ├── mall-order-application
│   ├── mall-order-domain
│   ├── mall-order-infrastructure
│   └── mall-order-adapter
├── mall-marketing -- 营销领域微服务
│   ├── mall-marketing-api
│   ├── mall-marketing-application
│   ├── mall-marketing-domain
│   ├── mall-marketing-infrastructure
│   └── mall-marketing-adapter
├── mall-content -- 内容领域微服务
│   ├── mall-content-api
│   ├── mall-content-application
│   ├── mall-content-domain
│   ├── mall-content-infrastructure
│   └── mall-content-adapter
├── mall-search -- 搜索领域微服务
│   ├── mall-search-api
│   ├── mall-search-application
│   ├── mall-search-domain
│   ├── mall-search-infrastructure
│   └── mall-search-adapter
└── mall-monitor -- 系统监控服务
```

## 技术架构

### 系统架构图

![系统架构图](../document/resource/mall_micro_service_arch.jpg)

### 业务架构图

![业务架构图](../document/resource/re_mall_business_arch.jpg)

### 技术栈

#### 微服务基础架构

- Spring Cloud Alibaba - 微服务框架
- Spring Cloud Gateway - API网关
- Nacos - 服务注册与配置中心
- Sentinel - 服务熔断与限流
- Seata - 分布式事务
- Spring Cloud Stream - 消息驱动
- Spring Cloud Sleuth + Zipkin - 分布式链路追踪
- Spring Boot Admin - 微服务管理与监控

#### 后端技术

- SpringBoot 2.7.5
- SpringSecurity + OAuth2
- MyBatis + MyBatisPlus
- Elasticsearch
- RocketMQ/RabbitMQ
- Redis + Redisson
- MongoDB
- Kafka
- ShardingSphere - 分库分表
- MinIO - 对象存储
- LogStash + ELK
- Docker + Kubernetes
- Prometheus + Grafana - 监控
- JWT
- MapStruct - DTO转换
- Lombok
- Hutool
- Swagger/Knife4j - API文档

#### 前端技术

- Vue3
- Vue-router
- Pinia/Vuex
- Element-Plus
- Axios
- ECharts
- Vite

#### 开发环境

- JDK 1.8+/JDK 17
- MySQL 5.7+/8.0
- Redis 7.0+
- MongoDB 5.0+
- RocketMQ/RabbitMQ
- Elasticsearch 7.17.3
- Logstash 7.17.3
- Kibana 7.17.3
- Nacos 2.1.0
- Sentinel 1.8.5
- Seata 1.5.2

## 环境准备

在开始构建项目之前，需要准备以下环境：

### 基础环境

1. **JDK 1.8+/JDK 17**

   - 下载地址：https://www.oracle.com/java/technologies/downloads/
   - 安装并配置环境变量
2. **Maven 3.8+**

   - 下载地址：https://maven.apache.org/download.cgi
   - 安装并配置环境变量
   - 建议配置阿里云Maven镜像加速下载：

   ```xml
   <mirror>
       <id>aliyunmaven</id>
       <mirrorOf>central</mirrorOf>
       <name>aliyun</name>
       <url>https://maven.aliyun.com/repository/public</url>
   </mirror>
   ```
3. **MySQL 5.7+/8.0**

   - 下载地址：https://www.mysql.com/downloads/
   - 安装并创建数据库mall
   - 导入项目中的/document/sql/mall.sql文件
4. **Redis 7.0+**

   - 下载地址：https://redis.io/download
   - 安装并启动Redis服务

### 微服务必需环境

5. **Nacos 2.1.0**

   - 下载地址：https://github.com/alibaba/nacos/releases
   - 配置并启动Nacos服务

   ```bash
   # 单机模式启动
   sh startup.sh -m standalone
   ```
6. **Sentinel 1.8.5**

   - 下载地址：https://github.com/alibaba/Sentinel/releases
   - 启动Sentinel控制台

   ```bash
   java -Dserver.port=8858 -jar sentinel-dashboard-1.8.5.jar
   ```
7. **Seata 1.5.2**

   - 下载地址：https://github.com/seata/seata/releases
   - 配置并启动Seata服务器

### 可选环境

8. **MongoDB 5.0+**

   - 下载地址：https://www.mongodb.com/download-center
   - 安装并启动MongoDB服务
9. **RocketMQ 4.9.4**

   - 下载地址：https://rocketmq.apache.org/download
   - 安装并启动RocketMQ服务

   ```bash
   # 启动NameServer
   nohup sh bin/mqnamesrv &
   # 启动Broker
   nohup sh bin/mqbroker -n localhost:9876 &
   ```
10. **Elasticsearch 7.17.3**

    - 下载地址：https://www.elastic.co/downloads/elasticsearch
    - 安装并启动Elasticsearch服务
11. **Logstash 7.17.3** (可选)

    - 下载地址：https://www.elastic.co/cn/downloads/logstash
    - 安装并配置Logstash
12. **Kibana 7.17.3** (可选)

    - 下载地址：https://www.elastic.co/cn/downloads/kibana
    - 安装并启动Kibana服务
13. **Zipkin** (可选)

    - 下载地址：https://zipkin.io/pages/quickstart.html
    - 启动Zipkin服务

    ```bash
    java -jar zipkin.jar
    ```
14. **MinIO** (可选)

    - 下载地址：https://min.io/download
    - 启动MinIO服务

    ```bash
    ./minio server /data
    ```

## 项目构建步骤

### 1. 克隆代码

```bash
git clone https://github.com/macrozheng/mall-microservice.git
```

### 2. 导入项目

- 使用IDEA或Eclipse导入项目
- 确保Maven正确配置

### 3. 配置Nacos

1. 启动Nacos服务
2. 导入项目配置到Nacos
   - 访问 http://localhost:8848/nacos
   - 创建mall命名空间
   - 导入各个微服务的配置文件

### 4. 配置数据库

- 根据SQL脚本创建各微服务所需的数据库
- 用户微服务：mall_ums
- 商品微服务：mall_pms
- 订单微服务：mall_oms
- 营销微服务：mall_sms
- 内容微服务：mall_cms

### 5. 构建项目

在项目根目录执行Maven命令：

```bash
mvn clean package -DskipTests
```

该命令会跳过测试并打包所有模块。

## 运行各个微服务

### 基础服务启动顺序

1. **启动Nacos**
2. **启动Sentinel**
3. **启动Seata**
4. **启动Zipkin**（可选）
5. **启动Redis**
6. **启动MySQL**

### 微服务启动顺序

1. **服务网关 (mall-gateway)**

   ```bash
   cd mall-infrastructure/mall-gateway
   mvn spring-boot:run
   ```
2. **认证服务 (mall-auth)**

   ```bash
   cd mall-infrastructure/mall-auth
   mvn spring-boot:run
   ```
3. **用户服务 (mall-user)**

   ```bash
   cd mall-user
   mvn spring-boot:run
   ```
4. **商品服务 (mall-product)**

   ```bash
   cd mall-product
   mvn spring-boot:run
   ```
5. **订单服务 (mall-order)**

   ```bash
   cd mall-order
   mvn spring-boot:run
   ```
6. **营销服务 (mall-marketing)**

   ```bash
   cd mall-marketing
   mvn spring-boot:run
   ```
7. **内容服务 (mall-content)**

   ```bash
   cd mall-content
   mvn spring-boot:run
   ```
8. **搜索服务 (mall-search)**

   ```bash
   cd mall-search
   mvn spring-boot:run
   ```

### 验证服务

1. 访问Nacos控制台查看服务注册情况：http://localhost:8848/nacos
2. 访问API网关：http://localhost:9201
3. API文档访问：

   - 通过网关访问Swagger：http://localhost:9201/doc.html
   - 各微服务单独文档：
     - 用户服务：http://localhost:8201/doc.html
     - 商品服务：http://localhost:8202/doc.html
     - 订单服务：http://localhost:8203/doc.html
     - 营销服务：http://localhost:8204/doc.html

## 微服务DDD架构详解

### 微服务内部结构

每个微服务采用六边形架构，内部结构划分如下：

#### API模块（mall-xxx-api）

- 定义服务对外暴露的接口
- 包含DTO（数据传输对象）
- 为其他微服务提供依赖

示例（商品服务API）：

```java
@FeignClient(value = "mall-product", path = "/product")
public interface ProductFeignClient {
  
    @GetMapping("/{id}")
    Result<ProductDTO> getProduct(@PathVariable("id") Long id);
  
    @GetMapping("/list")
    Result<CommonPage<ProductDTO>> list(@RequestParam(required = false) String keyword,
                                       @RequestParam(required = false) Long brandId,
                                       @RequestParam(required = false) Long productCategoryId,
                                       @RequestParam(required = false) Integer pageNum,
                                       @RequestParam(required = false) Integer pageSize,
                                       @RequestParam(required = false) String sort);
}
```

#### 领域模块（mall-xxx-domain）

- 包含领域实体（Entity）
- 值对象（Value Object）
- 领域服务（Domain Service）
- 领域事件（Domain Event）
- 仓储接口（Repository Interface）

示例（商品领域模型）：

```java
public class Product {
    // 标识
    private ProductId id;
    // 商品名称
    private ProductName name;
    // 品牌
    private BrandId brandId;
    // 分类
    private CategoryId categoryId;
    // 价格
    private Money price;
    // 库存
    private Stock stock;
    // 商品属性集合
    private Set<ProductAttribute> attributes;
  
    // 领域行为
    public void decreaseStock(int quantity) {
        this.stock = this.stock.decrease(quantity);
        DomainEventPublisher.publish(new ProductStockChangedEvent(this.id, quantity));
    }
  
    public void updatePrice(Money newPrice) {
        if (this.price.equals(newPrice)) {
            return;
        }
        Money oldPrice = this.price;
        this.price = newPrice;
        DomainEventPublisher.publish(new ProductPriceChangedEvent(this.id, oldPrice, newPrice));
    }
  
    // 工厂方法
    public static Product create(ProductId id, ProductName name, BrandId brandId, 
                                 CategoryId categoryId, Money price, Stock stock) {
        Product product = new Product();
        product.id = id;
        product.name = name;
        product.brandId = brandId;
        product.categoryId = categoryId;
        product.price = price;
        product.stock = stock;
        product.attributes = new HashSet<>();
      
        DomainEventPublisher.publish(new ProductCreatedEvent(product));
        return product;
    }
}
```

#### 应用模块（mall-xxx-application）

- 应用服务（Application Service）协调领域对象
- 事务边界
- 权限检查
- 领域事件处理
- 使用CQRS模式分离读写操作

示例（商品应用服务）：

```java
@Service
@Transactional
public class ProductApplicationService {
  
    private final ProductRepository productRepository;
    private final BrandRepository brandRepository;
    private final CategoryRepository categoryRepository;
    private final ProductDomainService productDomainService;
  
    public ProductApplicationService(ProductRepository productRepository,
                                     BrandRepository brandRepository,
                                     CategoryRepository categoryRepository,
                                     ProductDomainService productDomainService) {
        this.productRepository = productRepository;
        this.brandRepository = brandRepository;
        this.categoryRepository = categoryRepository;
        this.productDomainService = productDomainService;
    }
  
    public ProductId createProduct(CreateProductCommand command) {
        // 验证品牌和分类是否存在
        Brand brand = brandRepository.findById(command.getBrandId())
            .orElseThrow(() -> new EntityNotFoundException("品牌不存在"));
      
        Category category = categoryRepository.findById(command.getCategoryId())
            .orElseThrow(() -> new EntityNotFoundException("分类不存在"));
      
        // 创建商品领域对象
        Product product = Product.create(
            productRepository.nextId(),
            new ProductName(command.getName()),
            brand.getId(),
            category.getId(),
            new Money(command.getPrice()),
            new Stock(command.getStock())
        );
      
        // 添加商品属性
        for (ProductAttributeParam attributeParam : command.getAttributes()) {
            product.addAttribute(
                new ProductAttribute(attributeParam.getKey(), attributeParam.getValue())
            );
        }
      
        // 持久化
        productRepository.save(product);
      
        return product.getId();
    }
  
    public void updateProductPrice(UpdateProductPriceCommand command) {
        Product product = productRepository.findById(command.getProductId())
            .orElseThrow(() -> new EntityNotFoundException("商品不存在"));
      
        // 执行领域行为
        product.updatePrice(new Money(command.getNewPrice()));
      
        // 持久化
        productRepository.save(product);
    }
}
```

#### 基础设施模块（mall-xxx-infrastructure）

- 仓储实现（Repository Implementation）
- 消息中间件集成
- 缓存实现
- 外部服务集成
- 持久化映射

示例（商品仓储实现）：

```java
@Repository
public class ProductRepositoryImpl implements ProductRepository {
  
    private final ProductMapper productMapper;
    private final ProductAttributeMapper productAttributeMapper;
    private final IdGenerator idGenerator;
    private final ProductDOConverter converter;
  
    @Override
    public Optional<Product> findById(ProductId productId) {
        ProductDO productDO = productMapper.selectById(productId.getValue());
        if (productDO == null) {
            return Optional.empty();
        }
      
        List<ProductAttributeDO> attributes = 
            productAttributeMapper.selectByProductId(productId.getValue());
      
        return Optional.of(converter.toDomain(productDO, attributes));
    }
  
    @Override
    public void save(Product product) {
        ProductDO productDO = converter.toDO(product);
      
        if (productMapper.selectById(product.getId().getValue()) == null) {
            productMapper.insert(productDO);
        } else {
            productMapper.updateById(productDO);
        }
      
        // 保存属性
        saveAttributes(product);
    }
  
    @Override
    public ProductId nextId() {
        return new ProductId(idGenerator.nextId());
    }
  
    private void saveAttributes(Product product) {
        // 实现属性保存逻辑
    }
}
```

#### 适配器模块（mall-xxx-adapter）

- REST控制器（Controller）
- 消息监听器（MessageListener）
- 调度任务（Scheduler）
- 外部接口适配器

示例（商品控制器）：

```java
@RestController
@RequestMapping("/product")
@Api(tags = "商品管理")
public class ProductController {
  
    private final ProductApplicationService productApplicationService;
    private final ProductQueryService productQueryService;
  
    @PostMapping
    @ApiOperation("创建商品")
    public Result<Long> createProduct(@RequestBody @Valid CreateProductCommand command) {
        ProductId productId = productApplicationService.createProduct(command);
        return Result.success(productId.getValue());
    }
  
    @PutMapping("/{id}/price")
    @ApiOperation("更新商品价格")
    public Result<Void> updateProductPrice(@PathVariable Long id,
                                          @RequestBody @Valid UpdateProductPriceCommand command) {
        command.setProductId(id);
        productApplicationService.updateProductPrice(command);
        return Result.success();
    }
  
    @GetMapping("/{id}")
    @ApiOperation("获取商品详情")
    public Result<ProductDTO> getProduct(@PathVariable Long id) {
        ProductDTO productDTO = productQueryService.getProduct(id);
        return Result.success(productDTO);
    }
  
    @GetMapping("/list")
    @ApiOperation("分页查询商品")
    public Result<CommonPage<ProductDTO>> list(ProductQueryParam queryParam) {
        Page<ProductDTO> page = productQueryService.query(queryParam);
        return Result.success(CommonPage.restPage(page));
    }
}
```

### 领域事件与集成

在微服务架构中，领域事件用于实现服务间的异步通信和数据一致性：

```java
// 领域事件定义
public class ProductCreatedEvent implements DomainEvent {
    private final ProductId productId;
    private final Instant occurredTime;
  
    public ProductCreatedEvent(Product product) {
        this.productId = product.getId();
        this.occurredTime = Instant.now();
    }
  
    @Override
    public Instant occurredTime() {
        return occurredTime;
    }
}

// 事件发布
@Component
public class ProductEventPublisher {
  
    private final StreamBridge streamBridge;
  
    public void publishProductCreatedEvent(ProductCreatedEvent event) {
        ProductCreatedMessage message = new ProductCreatedMessage(
            event.getProductId().getValue(),
            event.getOccurredTime()
        );
        streamBridge.send("productCreated-out-0", message);
    }
}

// 事件订阅（在营销服务中）
@Component
public class ProductEventConsumer {
  
    private final PromotionApplicationService promotionService;
  
    @Bean
    public Consumer<ProductCreatedMessage> productCreated() {
        return message -> {
            log.info("Received product created event: {}", message);
            promotionService.handleNewProduct(message.getProductId());
        };
    }
}
```

## Docker部署

项目支持使用Docker进行部署，采用Docker Compose管理多个微服务。

### 使用Maven插件构建Docker镜像

1. 确保Docker服务已启动
2. 配置Docker远程访问

   - 修改pom.xml中的docker.host属性为你的Docker服务地址
3. 构建镜像

   ```bash
   mvn clean package docker:build -DskipTests
   ```

### 使用Docker Compose部署

项目提供了docker-compose配置：

- 基础环境部署：使用 `document/docker/docker-compose-env.yml`
- 微服务部署：使用 `document/docker/docker-compose-app.yml`

部署步骤：

1. 部署基础环境：

   ```bash
   docker-compose -f document/docker/docker-compose-env.yml up -d
   ```
2. 部署微服务：

   ```bash
   docker-compose -f document/docker/docker-compose-app.yml up -d
   ```

### Kubernetes部署 (可选)

针对生产环境，我们也提供了Kubernetes部署配置，详见项目中的k8s目录。

### 前端项目部署

前端项目基于微服务架构进行了适配：

1. 克隆前端项目：

   - 后台管理前端：https://github.com/macrozheng/mall-admin-web
   - 前台商城前端：https://github.com/macrozheng/mall-app-web
2. 安装依赖并运行：

   ```bash
   npm install
   npm run dev
   ```
3. 构建生产环境：

   ```bash
   npm run build
   ```

## 微服务开发指南

### 创建新微服务的步骤

1. 创建领域模型
   - 根据限界上下文定义领域实体和值对象
